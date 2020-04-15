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
    [Parameter(Mandatory=$false)][bool]$installInfo = $true
)

function main(){
    Clear-Host
    if (-not $sqlDir.ToString().EndsWith("\")){
        $sqlDir = $($sqlDir.ToString()+'\' )
    }
    if (-not $outDir.ToString().EndsWith("\")){
        $outDir = $($outDir.ToString()+'\' )
    }

    $startDate = $(get-date)
    $queryList = new-object System.Collections.ArrayList;
    $dbQueryList = new-object System.Collections.ArrayList;

    #First Verify ProductVersion and if the instance is responsible
    $queryPath = $("$sqlDir"+'mssqlserver.instanceVersion.sql'); 
    $outputFile = $("$outDir"+'instanceVersion.json'); 
    .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType

    if (Test-Path $($outputFile)){
        $jsonInstance = Get-Content -Path $outputFile |ConvertFrom-Json
        Write-Output $('productVersion :' + $($jsonInstance.productVersion)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
    }
    if (-not [String]::IsNullOrEmpty($jsonInstance)){
        $sqlmajor = $jsonInstance.productVersion.split(".")
        
        #Adding queries to execution list
        #Greather than or equal 2017
        if($sqlmajor[0] -ge '14'){ 
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
        #Less than 2017
        else {
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

        $queryList |ForEach-Object {
            $itemList = $_;
            Write-Output $($($itemList.order.ToString()) + ':'+ $($itemList.queryPath) +':'+ $($itemList.outputFile)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
            
            $queryPath = $($itemList.queryPath)
            $outputType = $($itemList.outputType);
            $outputFile = $($itemList.outputFile) +'.'+ $outputType
            .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType
            
            #Execute for each database
            if ($($itemList.foreachDB)) { 
                if (Test-Path $($outputFile)){
                    $dbList = Get-Content -Path $outputFile |ConvertFrom-Json
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
                                .\ercoleAgentMsSQLServer-Query.ps1 -instance $instance -dbName $($dbInfo.name) -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType -Parameters $parameters
                            }
                        }else{
                            Write-Output $('Alert:' + $($dbInfo.name) +'is '+ $($dbInfo.state_desc)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
                        }
                    }
                }
            }
        }

        if ($installInfo){
            #Installed Patches
            Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Sort-Object -Property DisplayName `
            | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
            | Where-Object {($_.DisplayName -like "Hotfix*SQL*") -or ($_.DisplayName -like "Service Pack*SQL*")} `
            | ConvertTo-Json|Out-File -FilePath $("$outDir"+'patches.json')

            
            #Last Patch installed
            $queryPath = $("$sqlDir"+'mssqlserver.psu-1.sql'); 
            $outputFile = $("$outDir"+'psu-aux.json'); 
            .\ercoleAgentMsSQLServer-Query.ps1 -QueryPath $queryPath -OutputFile $outputFile -OutputType $outputType
            if (Test-Path $($outputFile)){
                $lastPatch = Get-Content -Path $outputFile |ConvertFrom-Json
                Write-Output $('Last Patch :' + $($lastPatch.ProductUpdateReference)) |Out-File -FilePath $("$outDir"+'gather.log') -Append
            }
            if (-not [String]::IsNullOrEmpty($lastPatch)){
                Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall `
                | Get-ItemProperty | Sort-Object -Property DisplayName `
                | Select-Object -Property DisplayName, DisplayVersion, InstallDate `
                | Where-Object {($_.DisplayName -like $('*'+$lastPatch.ProductUpdateReference+'*'))} `
                | ConvertTo-Json|Out-File -FilePath $("$outDir"+'psu-1.json')
            }
            Remove-item -Path $outputFile
            
            #Features
            .\ercoleAgentMsSQLServer-GetFeatures.ps1 -majorVersion $($sqlmajor[0]) -outDir $outDir
        }


    }

   

    $endDate = $(get-date)
    $diffDate= New-TimeSpan -Start $startDate -End $endDate
    Write-Output $("Elapsed time " + $($diffDate.TotalSeconds) + " s") |Out-File -FilePath $("$outDir"+'gather.log') -Append
 }              

main