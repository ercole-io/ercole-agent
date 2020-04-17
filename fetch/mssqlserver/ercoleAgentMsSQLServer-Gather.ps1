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
    [Parameter(Mandatory=$false)][string]$instance ="DESKTOP-QN4D1T4\SQL2014",
    [Parameter(Mandatory=$false)][string]$sqlDir ="..\..\sql\mssqlserver\",
    [Parameter(Mandatory=$false)][string]$outDir ="..\..\Output\",
    [Parameter(Mandatory=$false)][ValidateSet("json")] [string]$outputType ="json",
    [Parameter(Mandatory=$false)][ValidateSet("file","object")] [string]$outputAs ="object",
    [Parameter(Mandatory=$false)][ValidateSet("all","dbmounted", "edition", "listDatabases","db", "backup_schedule", "dbstatus", "schema", "ts", "segment_advisor","patches","psu-1","sqlFeatures")][string]$action = "all"
)

function executeQueue([System.Collections.ArrayList]$queryList, [System.Collections.ArrayList]$dbQueryList){
    $resultList = new-object System.Collections.ArrayList;
    $queryList |ForEach-Object {
        $itemList = $_;
        Write-Output $($($itemList.order.ToString()) + ':'+ $($itemList.queryPath) +':'+ $($itemList.outputFile)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
        
        $queryPath = $($itemList.queryPath)
        $outputType = $($itemList.outputType);
        $outputFile = $($itemList.outputFile) +'.'+ $outputType

        #Execute for each database
        if (-not $($itemList.foreachDB)) {
            if ($outputAs -eq "file"){
                .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -outputAs $outputAs|ConvertFrom-Json
            }else{
                $execResult = .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -outputAs object|ConvertFrom-Json
                $auxLayout = $($outputFile)
                $auxLayout = $auxLayout.Substring($auxLayout.LastIndexOf('\')+1, $auxLayout.Length - $auxLayout.LastIndexOf('\')-1)
                $jsonbase = @{layout=$auxLayout}
                $jsonbase.Add("data",$execResult)
                $jsonbase| ConvertTo-Json -Depth 10
    }
        }else{
            $dbList = .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -outputAs object|ConvertFrom-Json
            if (-not [string]::IsNullOrEmpty($dbList)){
                $dbList |ForEach-Object{
                    $dbInfo = $_
                    Write-Output $('dbname:' + $($dbInfo.name)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
                    if ($($dbInfo.state_desc) -eq 'ONLINE'){
                        $dbQueryList| ForEach-Object{
                            $dbItemList = $_;
                            $queryPath = $($dbItemList.queryPath);
                            $outputType = $($dbItemList.outputType); 
                            $outputFile = $( $($dbItemList.outputFile) +'_'+ $($dbInfo.name) +'.'+ $outputType );
                            
                            $parameters = @{
                                pdatabase_id =  @{name='pdatabase_id'; value=$($dbInfo.database_id); datatype='Int'; size=$null; precision=$null}
                            }	
                            if ($outputAs -eq "file"){
                                .\ercoleAgentMsSQLServer-Query.ps1 -instance $instance -dbName $($dbInfo.name) -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -Parameters $parameters -outputAs $outputAs
                            }else{
                                
                                $execResult = .\ercoleAgentMsSQLServer-Query.ps1 -instance $instance -dbName $($dbInfo.name) -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -Parameters $parameters -outputAs object|ConvertFrom-Json
                                $auxLayout = $($dbItemList.outputFile)
                                $auxLayout = $auxLayout.Substring($auxLayout.LastIndexOf('\')+1, $auxLayout.Length - $auxLayout.LastIndexOf('\')-1)
                                $jsonbase = @{layout=$auxLayout; database_name=$($dbInfo.name)}
                                $jsonbase.Add("data",$execResult)
                                $resultList.Add($jsonbase) |Out-Null
                            }
                     }
                    }else{
                        Write-Output $('Alert:' + $($dbInfo.name) +'is '+ $($dbInfo.state_desc)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
                    }
                }
                if ($outputAs -eq "object"){
                    #$jsonbase = @{layout=;data=$outputList}
                    #$jsonbase| ConvertTo-Json -Depth 10
                    $resultList| ConvertTo-Json -Depth 10
                }
            }
        }
    }
}


 function installedPatches(){
    #Installed Patches
    $jsonbase = @{layout='patches'}
    $execResult = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Sort-Object -Property DisplayName `
        | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
        | Where-Object {($_.DisplayName -like "Hotfix*SQL*") -or ($_.DisplayName -like "Service Pack*SQL*")} 

    $jsonbase.Add("data",$execResult)                            
    if ($outputAs -eq "file"){
        $jsonbase | ConvertTo-Json -Depth 10|Out-File -FilePath $("$outDir"+'patches.json')
    }else{
        $jsonbase| ConvertTo-Json -Depth 10
    }
}


function lastInstalledPatch(){
    #Last Patch installed
    $jsonbase = @{layout='psu-1'}
    
    $queryPath = $("$sqlDir"+'mssqlserver.psu-1.sql'); 
    $outputFile = $("$outDir"+'psu-aux.json'); 
    $lastPatch = .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -outputAs object|ConvertFrom-Json
    Write-Output $('Last Patch :' + $($lastPatch.ProductUpdateReference)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
    if (-not [String]::IsNullOrEmpty($lastPatch)){
        $execResult = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall `
            | Get-ItemProperty | Sort-Object -Property DisplayName `
            | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
            | Where-Object {($_.DisplayName -like $('*'+$lastPatch.ProductUpdateReference+'*'))}

        $jsonbase.Add("data",$execResult)
        if ($outputAs -eq "file"){
            $jsonbase| ConvertTo-Json -Depth 10|Out-File -FilePath $("$outDir"+'psu-1.json')
        }else{
            $jsonbase| ConvertTo-Json -Depth 10
        }
    }
}


function sqlFeatures() {
    #Features
    $jsonbase = @{layout='sqlFeatures'}
    $execResult = .\ercoleAgentMsSQLServer-GetFeatures.ps1 -majorVersion $($sqlmajor[0]) -outDir $outDir -outputAs object|ConvertFrom-Json
    $jsonbase.Add("data",$execResult)
        if ($outputAs -eq "file"){
            $jsonbase| ConvertTo-Json -Depth 10|Out-File -FilePath $("$outDir"+'sqlFeatures.json')
        }else{
            $jsonbase| ConvertTo-Json -Depth 10
        }
}
            
       

 #"all","dbmounted", "edition", "listDatabases","db", "backup_schedule", "dbstatus", "schema", "ts", "segment_advisor", "patches", "psu-1"
 function main(){
    Clear-Host
    $startDate = $(get-date)
    Write-Output $(">>Start at: "+$startDate.ToString('yyyyMMdd-HH:mm:ss'))|Out-File -FilePath $("$outDir"+'gather.log') -Append
    
    if (-not $sqlDir.ToString().EndsWith("\")){
        $sqlDir = $($sqlDir.ToString()+'\' )
    }
    if (-not $outDir.ToString().EndsWith("\")){
        $outDir = $($outDir.ToString()+'\' )
    }
    $queryList = new-object System.Collections.ArrayList;
    $dbQueryList = new-object System.Collections.ArrayList;
    $installedPatches = $false
    $lastInstalledPatch = $false
    $sqlFeatures = $false
    
    #First Verify ProductVersion and if the instance is up
    $queryPath = $("$sqlDir"+'mssqlserver.instanceVersion.sql'); 
    $outputFile = $("$outDir"+'instanceVersion.json'); 
    $jsonInstance = .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -outputAs object|ConvertFrom-Json
    Write-Output $('productVersion :' + $($jsonInstance.productVersion)) |Out-File -FilePath $("$outDir"+'gather.log') -Append

    #Adding queries to execution list
    if (-not [String]::IsNullOrEmpty($jsonInstance)){
        $sqlmajor = $jsonInstance.productVersion.split(".")
        if ($action -eq "all"){
            if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.14.sql'); outputFile= $("$outDir"+'dbmounted'); outputType='json'; foreachDB=$false}) |Out-Null
                $queryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); outputFile= $("$outDir"+'edition'); outputType='json'; foreachDB=$false}) |Out-Null
                $queryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
                            
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.14.sql'); outputFile= $("$outDir"+'db'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); outputFile= $("$outDir"+'backup_schedule'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); outputFile= $("$outDir"+'dbstatus'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); outputFile= $("$outDir"+'schema'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=5; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); outputFile= $("$outDir"+'ts'); outputType='json'}) |Out-Null
		    	$dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); outputFile= $("$outDir"+'segment_advisor'); outputType='json'}) |Out-Null
            }
            else { #older than 2017
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.10.sql'); outputFile= $("$outDir"+'dbmounted'); outputType='json'; foreachDB=$false}) |Out-Null
                $queryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); outputFile= $("$outDir"+'edition'); outputType='json'; foreachDB=$false}) |Out-Null
                $queryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
                
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.10.sql'); outputFile= $("$outDir"+'db'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=2; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); outputFile= $("$outDir"+'backup_schedule'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=3; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); outputFile= $("$outDir"+'dbstatus'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=4; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); outputFile= $("$outDir"+'schema'); outputType='json'}) |Out-Null
                $dbQueryList.Add(@{order=5; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); outputFile= $("$outDir"+'ts'); outputType='json'}) |Out-Null
		    	$dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); outputFile= $("$outDir"+'segment_advisor'); outputType='json'}) |Out-Null
            }
            $installedPatches = $true
            $lastInstalledPatch = $true
            $sqlFeatures = $true

        }elseif ($action -eq "dbmounted"){
            if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.14.sql'); outputFile= $("$outDir"+'dbmounted'); outputType='json'; foreachDB=$false}) |Out-Null
            }
            else { #older than 2017
                $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbmounted.10.sql'); outputFile= $("$outDir"+'dbmounted'); outputType='json'; foreachDB=$false}) |Out-Null
            }

        }elseif ($action -eq "edition"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.edition.sql'); outputFile= $("$outDir"+'edition'); outputType='json'; foreachDB=$false}) |Out-Null

        }elseif ($action -eq "listDatabases"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$false})|Out-Null

        }elseif ($action -eq "db"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            if($sqlmajor[0] -ge '14'){ #newer or equal than 2017
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.14.sql'); outputFile= $("$outDir"+'db'); outputType='json'}) |Out-Null
            }
            else { #older than 2017
                $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.db.10.sql'); outputFile= $("$outDir"+'db'); outputType='json'}) |Out-Null
            }

        }elseif ($action -eq "backup_schedule"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.backup_schedule.sql'); outputFile= $("$outDir"+'backup_schedule'); outputType='json'}) |Out-Null

        }elseif ($action -eq "dbstatus"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.dbstatus.sql'); outputFile= $("$outDir"+'dbstatus'); outputType='json'}) |Out-Null

        }elseif ($action -eq "schema"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.schema.sql'); outputFile= $("$outDir"+'schema'); outputType='json'}) |Out-Null

        }elseif ($action -eq "ts"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            $dbQueryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.ts.sql'); outputFile= $("$outDir"+'ts'); outputType='json'}) |Out-Null
		    	
        }elseif ($action -eq "segment_advisor"){
            $queryList.Add(@{order=1; queryPath=$("$sqlDir"+'mssqlserver.listDatabases.sql'); outputFile= $("$outDir"+'listDatabases'); outputType='json'; foreachDB=$true})|Out-Null
            $dbQueryList.Add(@{order=6; queryPath=$("$sqlDir"+'mssqlserver.segment_advisor.sql'); outputFile= $("$outDir"+'segment_advisor'); outputType='json'}) |Out-Null

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

    $endDate = $(get-date)
    Write-Output $(">>End at: "+$endDate.ToString('yyyyMMdd-HH:mm:ss'))|Out-File -FilePath $("$outDir"+'gather.log') -Append
    $diffDate= New-TimeSpan -Start $startDate -End $endDate
    Write-Output $(">>Elapsed time " + $($diffDate.TotalSeconds) + " s") |Out-File -FilePath $("$outDir"+'gather.log') -Append
 }
main