# Copyright (c) 2020 Sorint.lab S.p.A.
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

param(
    [Parameter(Mandatory=$false)][string]$instance = $null,
    [Parameter(Mandatory=$false)][string]$sqlDir =".\sql\mssqlserver",
    [Parameter(Mandatory=$false)][ValidateSet("listInstances", "all","dbmounted", "edition", "licensingInfo", "listDatabases","db", "backup_schedule", "dbstatus", "schema", "ts", "segment_advisor","patches","psu-1","sqlFeatures")][string]$action = "listInstances"
)


#Query functions################################################################
function SQLQuery{ 
    [CmdletBinding()] 
    param( 
        [Parameter(Position=0, Mandatory=$true)] [string]$ServerInstance, 
        [Parameter(Position=1, Mandatory=$false)] [string]$Database, 
        [Parameter(Position=2, Mandatory=$false)] [string]$Query, 
        [Parameter(Position=3, Mandatory=$false)] [string]$Username, 
        [Parameter(Position=4, Mandatory=$false)] [string]$Password, 
        [Parameter(Position=5, Mandatory=$false)] [Int32]$QueryTimeout=0, 
        [Parameter(Position=6, Mandatory=$false)] [Int32]$ConnectionTimeout=3600, 
        [Parameter(Position=7, Mandatory=$false)] [ValidateScript({test-path $_})] [string]$InputFile, 
        [Parameter(Position=8, Mandatory=$false)] [ValidateSet("DataSet", "DataTable", "DataRow")] [string]$As="DataRow" ,
        [Parameter(Position=9, Mandatory=$false)] [Object]$Parameters 
    ) 
 
    $conn=new-object System.Data.SqlClient.SQLConnection 
   
    if ($Username){
        $ConnectionString = "Server={0};Database={1};User ID={2};Password={3};Trusted_Connection=False;Connect Timeout={4}" -f $ServerInstance,$Database,$Username,$Password,$ConnectionTimeout 
    } else { 
        $ConnectionString = "Server={0};Database={1};Integrated Security=True;Connect Timeout={2}" -f $ServerInstance,$Database,$ConnectionTimeout 
    } 
 
    $conn.ConnectionString=$ConnectionString 
   
    if ($PSBoundParameters.Verbose) { 
        $conn.FireInfoMessageEventOnUserErrors=$true 
        $handler = [System.Data.SqlClient.SqlInfoMessageEventHandler] {Write-Verbose "$($_)"} 
        $conn.add_InfoMessage($handler) 
    } 
   
    try {
        $conn.Open()
        $cmd=new-object system.Data.SqlClient.SqlCommand 
        $cmd.Connection = $conn
        $cmd.CommandText = $Query
        $cmd.CommandTimeout=$QueryTimeout 

        if ($Parameters){
            foreach($par in $Parameters.Keys){
                $sqlParameter = New-Object System.Data.SqlClient.SqlParameter("@$par",$('[Data.SQLDBType]::'+$($Parameters[$par].datatype)))
                if (-Not([string]::IsNullOrEmpty($Parameters[$par].size))){
                    $sqlParameter.Size = $($Parameters[$par].size)
                }elseif (-Not([string]::IsNullOrEmpty($Parameters[$par].precision))){
                    $sqlParameter.Precision = $($Parameters[$par].precision)
                }
                $sqlParameter.Value = $($Parameters[$par].value)
                $cmd.Parameters.Add($sqlParameter) |Out-Null
            }
        }
    
        $ds=New-Object system.Data.DataSet 
        $da=New-Object system.Data.SqlClient.SqlDataAdapter($cmd) 
        [void]$da.fill($ds)
        $conn.Close() 
        switch ($As) 
        { 
          'DataSet'  { Write-Output ($ds) } 
          'DataTable' { Write-Output ($ds.Tables) } 
          'DataRow'  { Write-Output ($ds.Tables[0]) } 
        }
    } catch{ 
        $ex = $_.Exception 
        Write-Error "$ex.Message" 
    } 
} 

function getQuery(
    [Parameter(Mandatory=$false)][string]$dbName ="master",
    [Parameter(Mandatory=$false)][string]$username = $null,
    [Parameter(Mandatory=$false)][string]$passwordEncryp = $null,
    [Parameter(Mandatory=$false)][string]$queryPath = "",
    [Parameter(Mandatory=$false)] [Object]$parameters =$null
){
    if (Test-Path $queryPath){
        $query = Get-Content $queryPath |Out-String
        $results = SQLQuery -ServerInstance $instance -Database $dbName -Query $query -as DataTable -Username $username -Password $passwordEncryp -Parameters $Parameters
        $results |Select $results.Columns.ColumnName |ConvertTo-Json
    }
 }              

#SQLFeatures functions##########################################################
function convertReportToJson([string]$setupBootstrap){
    #$SetupBootstrap ='C:\Program Files\Microsoft SQL Server\120\Setup Bootstrap'
    $SetupBootstrapLog = [System.IO.Path]::Combine($SetupBootstrap, 'Log')
    $lastDir = Get-ChildItem -LiteralPath $SetupBootstrapLog -Directory |Sort LastWriteTime |Select name -last 1
    $reportPath = [System.IO.Path]::Combine($SetupBootstrapLog, $lastDir.Name)
    $reportFile = $($reportPath+'\SqlDiscoveryReport.xml')
    $result = $null
    if ([System.IO.File]::Exists($reportFile)){
        [xml]$xvar = Get-Content $reportFile
        $result = $xvar.ArrayOfDiscoveryInformation.DiscoveryInformation `
            |Select Product, Instance, InstanceID, Feature, Language, Edition, Version, Clustered, Configured
    }
    return $result
}

function getNewerVersionPath([string]$targetVersion){
    $MSSQLpath = [System.IO.Path]::Combine($env:ProgramFiles, "Microsoft SQL Server")
    $lstCompatLevelDirs = New-Object "System.Collections.Generic.List[Int32]"
    
    Get-ChildItem -Directory $MSSQLpath -Filter $($targetVersion.ToString()+'*') | 
        ForEach-Object {
            [Int32]$DirNum = 0
    
            if ([Int32]::TryParse($_.Name, [ref]$DirNum))
            {
                $lstCompatLevelDirs.Add($DirNum)
            }
        }
    
    $lstCompatLevelDirs.Sort()
    $lstCompatLevelDirs.Reverse()
    
    [bool] $setupExeFound = $false
    
    ForEach($int in $lstCompatLevelDirs)
    {
        [string]$SetupBootstrap = [System.IO.Path]::Combine(
            [System.IO.Path]::Combine($MSSQLpath, $int.ToString()),
            "Setup Bootstrap")
    
        if ([System.IO.Directory]::Exists($SetupBootstrap))
        {
            ForEach($sqlSubDir in [System.IO.Directory]::GetDirectories($SetupBootstrap, "SQL*"))
            {
                [string]$setupExe = [System.IO.Path]::Combine($sqlSubDir, "setup.exe")
    
                if ([System.IO.File]::Exists($setupExe))
                {
                    $setupExeFound = $true
                    Start-Process -FilePath $setupExe -ArgumentList "/q /Action=RunDiscovery" -Wait
                    $finalResult = convertReportToJson -setupBootstrap $SetupBootstrap
                    break
                }
            }
        }
    
        if($setupExeFound)
        {
            break
        }
    }
    return $finalResult
}


function getOlderVersionPath([string]$targetVersion){
    $MSSQLpath = [System.IO.Path]::Combine($env:ProgramFiles, "Microsoft SQL Server")
    $lstCompatLevelDirs = New-Object "System.Collections.Generic.List[Int32]"
    [bool] $setupExeFound = $false
    if(-Not $setupExeFound)
    {
        $lstOldSqlVersionSetupExePaths = New-Object "System.Collections.Generic.List[string]"
    
        #SQL 2008
        $lstOldSqlVersionSetupExePaths.Add([System.IO.Path]::Combine($MSSQLpath, "100\Setup Bootstrap\Release\Setup.exe"))
    
        #SQL 2005
        $lstOldSqlVersionSetupExePaths.Add([System.IO.Path]::Combine($MSSQLpath, "90\Setup Bootstrap\Setup.exe"))
    
        ForEach($setupExe in $lstOldSqlVersionSetupExePaths)
        {
            if ([System.IO.File]::Exists($setupExe))
            {
                $setupExeFound = $true
                Start-Process -FilePath $setupExe -ArgumentList "/q /Action=RunDiscovery"
                break
            }
        }
    }
    return $setupExeFound
}


function getSQLFeatures(
    [Parameter(Mandatory=$false)][string]$majorVersion = ''
){
    $result = $null

    if ([string]::IsNullOrEmpty($majorVersion)){
        $result = getNewerVersionPath -targetVersion $null
        if (-not [string]::IsNullOrEmpty($result)){
            $result = getOlderVersionPath -targetVersion $null    
        }
    }elseif ($majorVersion -ge '10'){
        $result = getNewerVersionPath -targetVersion $majorVersion
    }else{
        $result = getOlderVersionPath -targetVersion $majorVersion
    }

    if (-not [string]::IsNullOrEmpty($result)){
        $result|ConvertTo-Json
    }
}



#Gather functions###############################################################
function executeQueue([System.Collections.ArrayList]$queryList, [System.Collections.ArrayList]$dbQueryList){
    $resultList = new-object System.Collections.ArrayList;
    $queryList |ForEach-Object {
        $itemList = $_;
       
        $queryPath = $($itemList.queryPath)
        $auxLayout = $($itemList.layout)
        
        #Execute for each database
        if (-not $($itemList.foreachDB)) {
            $execResult = getQuery -QueryPath $queryPath |ConvertFrom-Json
            $jsonbase = @{layout=$auxLayout}
            $jsonbase.Add("data",$execResult)
            $jsonbase| ConvertTo-Json -Depth 10
        }else{
            $dbList = getQuery -QueryPath $queryPath |ConvertFrom-Json
            if (-not [string]::IsNullOrEmpty($dbList)){
                $dbList |ForEach-Object{
                    $dbInfo = $_
                    if ($($dbInfo.state_desc) -eq 'ONLINE'){
                        $dbQueryList| ForEach-Object{
                            $dbItemList = $_;
                            $auxLayout = $($dbItemList.layout)
                            $queryPath = $($dbItemList.queryPath);
                            $parameters = @{
                                pdatabase_id =  @{name='pdatabase_id'; value=$($dbInfo.database_id); datatype='Int'; size=$null; precision=$null}
                            }	
                            $execResult = getQuery -dbName $($dbInfo.name) -QueryPath $queryPath -Parameters $parameters |ConvertFrom-Json
                            $jsonbase = @{layout=$auxLayout; database_name=$($dbInfo.name)}
                            $jsonbase.Add("data",$execResult)
                            $resultList.Add($jsonbase) |Out-Null
                        }
                    }
                }
                #Output
                $resultList| ConvertTo-Json -Depth 10
            }
        }
    }
}


 function installedPatches(){
    $jsonbase = @{layout='patches'}
    $execResult = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Sort-Object -Property DisplayName `
        | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
        | Where-Object {($_.DisplayName -like "Hotfix*SQL*") -or ($_.DisplayName -like "Service Pack*SQL*")} 

    $jsonbase.Add("data",$execResult)                            
    $jsonbase| ConvertTo-Json -Depth 10
}


function lastInstalledPatch(){
    $jsonbase = @{layout='psu-1'}
    $queryPath = $("$sqlDir"+'mssqlserver.psu-1.sql'); 
    $lastPatch = getQuery -QueryPath $queryPath |ConvertFrom-Json
    if (-not [String]::IsNullOrEmpty($lastPatch)){
        $execResult = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall `
            | Get-ItemProperty | Sort-Object -Property DisplayName `
            | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
            | Where-Object {($_.DisplayName -like $('*'+$lastPatch.ProductUpdateReference+'*'))}

        $jsonbase.Add("data",$execResult)
        $jsonbase| ConvertTo-Json -Depth 10
    }
}


function sqlFeatures() {
    $jsonbase = @{layout='sqlFeatures'}
    $execResult = getSQLFeatures -majorVersion $($sqlmajor[0]) |ConvertFrom-Json
    $jsonbase.Add("data",$execResult)
    $jsonbase| ConvertTo-Json -Depth 10
}


function listInstances(){
    $resultList = new-object System.Collections.ArrayList;
    $jsonbase = @{layout='listInstances'}
    $execResult = Get-Service |Select-Object Status, Name, DisplayName| where {$_.DisplayName -like 'SQL Server (*'}
    $hostname = $(hostname)
    $execResult|ForEach-Object{
        $currItem = $_
        if ($execResult.Name -eq 'MSSQLSERVER'){
            $connString = $hostname
        }else{
            $connString = $($currItem.Name.replace('MSSQL$',"$hostname\"))
        }
        if ($execResult.Status -eq 1){
            $statusDesc = 'Stopped'
        }elseif ($execResult.Status -eq 4){
            $statusDesc = 'Running'
        }else{
            $statusDesc = 'other'
        }
        $resultList.Add(@{name=$($currItem.Name); status=$statusDesc; displayName = $($currItem.DisplayName); connString=$connString})|Out-Null
    }
    
    $jsonbase.Add("data",$resultList)
    $jsonbase| ConvertTo-Json -Depth 10
}
            
       

function main(){
    Clear-Host
     
    if (-not $sqlDir.ToString().EndsWith("\")){
        $sqlDir = $($sqlDir.ToString()+'\' )
    }

    $queryList = new-object System.Collections.ArrayList;
    $dbQueryList = new-object System.Collections.ArrayList;
    $installedPatches = $false
    $lastInstalledPatch = $false
    $sqlFeatures = $false

    if ($action -eq "listInstances"){
        listInstances
    }else{
    
        #First Verify ProductVersion and if the instance is up
        $queryPath = $("$sqlDir"+'mssqlserver.instanceVersion.sql'); 
        $jsonInstance = getQuery -QueryPath $queryPath |ConvertFrom-Json

        #Adding queries to execution list
        if (-not [String]::IsNullOrEmpty($jsonInstance)){
            $sqlmajor = $jsonInstance.productVersion.split(".")
            if ($action -eq "all"){
                if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                    $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.14.sql'); layout= 'dbmounted';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); layout='edition';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.licensingInfo.sql'); layout='licensingInfo';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                                
                    $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.14.sql'); layout='db'}) |Out-Null
                    $dbQueryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); layout='backup_schedule'}) |Out-Null
                    $dbQueryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); layout='dbstatus'}) |Out-Null
                    $dbQueryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); layout='schema'}) |Out-Null
                    $dbQueryList.Add(@{order=5; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); layout='ts'}) |Out-Null
	    	    	$dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); layout='segment_advisor'}) |Out-Null
                }
                else { #older than 2017
                    $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.10.sql'); layout='dbmounted';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); layout='edition';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.licensingInfo.sql'); layout='licensingInfo';  foreachDB=$false}) |Out-Null
                    $queryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                    
                    $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.10.sql'); layout='db'}) |Out-Null
                    $dbQueryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); layout='backup_schedule'}) |Out-Null
                    $dbQueryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); layout='dbstatus'}) |Out-Null
                    $dbQueryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); layout='schema'}) |Out-Null
                    $dbQueryList.Add(@{order=5; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); layout='ts'}) |Out-Null
	    	    	$dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); layout='segment_advisor'}) |Out-Null
                }
                $installedPatches = $true
                $lastInstalledPatch = $true
                $sqlFeatures = $true

            }elseif ($action -eq "dbmounted"){
                if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                    $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.14.sql'); layout='dbmounted';  foreachDB=$false}) |Out-Null
                }
                else { #older than 2017
                    $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.10.sql'); layout='dbmounted';  foreachDB=$false}) |Out-Null
                }

            }elseif ($action -eq "edition"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); layout='edition';  foreachDB=$false}) |Out-Null

            }elseif ($action -eq "licensingInfo"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.licensingInfo.sql'); layout='licensingInfo';  foreachDB=$false}) |Out-Null

            }elseif ($action -eq "listDatabases"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$false})|Out-Null

            }elseif ($action -eq "db"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                    $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.14.sql'); layout='db'}) |Out-Null
                }
                else { #older than 2017
                    $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.10.sql'); layout='db'}) |Out-Null
                }

            }elseif ($action -eq "backup_schedule"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); layout='backup_schedule'}) |Out-Null

            }elseif ($action -eq "dbstatus"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); layout='dbstatus'}) |Out-Null

            }elseif ($action -eq "schema"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); layout='schema'}) |Out-Null

            }elseif ($action -eq "ts"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); layout='ts'}) |Out-Null
	    	    	
            }elseif ($action -eq "segment_advisor"){
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); layout='listDatabases';  foreachDB=$true})|Out-Null
                $dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); layout='segment_advisor'}) |Out-Null

            }elseif ($action -eq "patches"){
                $installedPatches = $true

            }elseif ($action -eq "psu-1"){
                $lastInstalledPatch = $true

            }elseif ($action -eq "sqlFeatures"){
                $sqlFeatures = $true

            }

            #executeQueue
            executeQueue -queryList $queryList -dbQueryList $dbQueryList
            if ($installedPatches){
                installedPatches
            }
            if ($lastInstalledPatch){
                lastInstalledPatch
            }
            if ($sqlFeatures){
                sqlFeatures
            }
        }
    }
 }
main