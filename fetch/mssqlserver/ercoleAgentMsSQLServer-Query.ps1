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
    [Parameter(Mandatory=$false)][string]$dbName ="master",
    [Parameter(Mandatory=$false)][string]$username = $null,
    [Parameter(Mandatory=$false)][string]$passwordEncryp = $null,
    [Parameter(Mandatory=$false)][string]$queryPath = "C:\Danilo\Sorint\Ercole-Agent\sql\mssqlserver\mssqlserver.dbmounted.10.sql",
    [Parameter(Mandatory=$false)][string]$outputFile = "C:\Danilo\Sorint\Ercole-Agent\Output\dbmounted.json",
    [Parameter(Mandatory=$false)][ValidateSet("csv", "json")] [string]$outputType ="json",
    [Parameter(Mandatory=$false)] [Object]$parameters =$null
)


function SQLQuery{ 
    [CmdletBinding()] 
    param( 
        [Parameter(Position=0, Mandatory=$true)] [string]$ServerInstance, 
        [Parameter(Position=1, Mandatory=$false)] [string]$Database, 
        [Parameter(Position=2, Mandatory=$false)] [string]$Query, 
        [Parameter(Position=3, Mandatory=$false)] [string]$Username, 
        [Parameter(Position=4, Mandatory=$false)] [string]$Password, 
        [Parameter(Position=5, Mandatory=$false)] [Int32]$QueryTimeout=0, 
        [Parameter(Position=6, Mandatory=$false)] [Int32]$ConnectionTimeout=15, 
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
   
    #Following EventHandler is used for PRINT and RAISERROR T-SQL statements. Executed when -Verbose parameter specified by caller 
    if ($PSBoundParameters.Verbose) { 
        $conn.FireInfoMessageEventOnUserErrors=$true 
        $handler = [System.Data.SqlClient.SqlInfoMessageEventHandler] {Write-Verbose "$($_)"} 
        $conn.add_InfoMessage($handler) 
    } 
   
    try {
        $conn.Open()
        $cmd=new-object system.Data.SqlClient.SqlCommand 
        $cmd.Connection = $conn
        #$cmd.CommandType = $CommandType
        $cmd.CommandText = $Query
        $cmd.CommandTimeout=$QueryTimeout 

        #newParameter
        if ($Parameters){
            foreach($par in $Parameters.Keys){
                #Write-host $($Parameters[$par].name)
                #Write-host $($Parameters[$par].datatype)
                #Write-host $($Parameters[$par].size)
                #Write-host $($Parameters[$par].precision)
                #Write-host $($Parameters[$par].value)
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
        #end newParameter

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
        #continue 
    } 
} 

function main(){
    #Clear-Host
    if (Test-Path $queryPath){
        $query = Get-Content $queryPath |Out-String

        $results = SQLQuery -ServerInstance $instance -Database $dbName -Query $query -as DataTable -Username $username -Password $passwordEncryp -Parameters $Parameters
        if ($OutputType -eq "json"){
            $results |Select $results.Columns.ColumnName |ConvertTo-Json |Out-File -FilePath $OutputFile
        }elseif ($OutputType -eq "csv"){
            $results|ConvertTo-Csv -NoTypeInformation -Delimiter ';'|Out-File -FilePath $OutputFile
        }

        
    }else{
        Write-Host "Alert: $QueryPath not found" -ForegroundColor Yellow
    }
 }              

main