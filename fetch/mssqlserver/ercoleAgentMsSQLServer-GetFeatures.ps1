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
    [Parameter(Mandatory=$false)][string]$majorVersion = '12',
    [Parameter(Mandatory=$false)][string]$outDir ="..\..\Output\",
    [Parameter(Mandatory=$false)][ValidateSet("file","object")] [string]$outputAs ="objetc"
)


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
        #$xvar.ArrayOfDiscoveryInformation.DiscoveryInformation `
        #    |Select Product, Instance, InstanceID, Feature, Language, Edition, Version, Clustered, Configured `
        #    |ConvertTo-Json|Out-File -LiteralPath  $("$outDir"+'hostSqlFeatures.json')
    }
    return $result
}

function getNewerVersionPath([string]$targetVersion){
    #Locate the "%PROGRAMFILES%\Microsoft SQL Server" folder.
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
    
    #Sort() the List, then Reverse() it so there is DESCENDING order.
    $lstCompatLevelDirs.Sort()
    $lstCompatLevelDirs.Reverse()
    
    [bool] $setupExeFound = $false
    
    <#
        Find the Setup Bootstrap Setup.exe file in the "highest" sub-folder.
        Here are a few examples:
            "%PROGRAMFILES%\Microsoft SQL Server\140\Setup Bootstrap\SQL2017\setup.exe"
            "%PROGRAMFILES%\Microsoft SQL Server\130\Setup Bootstrap\SQLServer2016\setup.exe"
            "%PROGRAMFILES%\Microsoft SQL Server\120\Setup Bootstrap\SQLServer2014\setup.exe"
            "%PROGRAMFILES%\Microsoft SQL Server\110\Setup Bootstrap\SQLServer2012\setup.exe"
            "%PROGRAMFILES%\Microsoft SQL Server\100\Setup Bootstrap\SQLServer2008R2\Setup.exe"
    #>
    ForEach($int in $lstCompatLevelDirs)
    {
        #The "Setup Bootstrap" path. For example: "%PROGRAMFILES%\Microsoft SQL Server\140\Setup Bootstrap
        [string]$SetupBootstrap = [System.IO.Path]::Combine(
            [System.IO.Path]::Combine($MSSQLpath, $int.ToString()),
            "Setup Bootstrap")
    
        if ([System.IO.Directory]::Exists($SetupBootstrap))
        {
            <#
                Iterate through the list of sub-folders with names that match the pattern "SQL*"
            #>
            ForEach($sqlSubDir in [System.IO.Directory]::GetDirectories($SetupBootstrap, "SQL*"))
            {
                <#
                    Search for "setup.exe". 
                    If found:
                        Run the exe with the appropriate parameters to run the discovery report.
                        Break out of the loops.
                #>
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
    #Locate the "%PROGRAMFILES%\Microsoft SQL Server" folder.
    $MSSQLpath = [System.IO.Path]::Combine($env:ProgramFiles, "Microsoft SQL Server")
    $lstCompatLevelDirs = New-Object "System.Collections.Generic.List[Int32]"
    [bool] $setupExeFound = $false
    <#
        If the Setup.exe is still not found, search for it in hard-coded paths that correspond
        to older versions that didn't use the current version/folder/naming convention.
    
        2008: "%PROGRAMFILES%\Microsoft SQL Server\100\Setup Bootstrap\Release\Setup.exe"
        2005: "%PROGRAMFILES%\Microsoft SQL Server\90\Setup Bootstrap\Setup.exe"
    #>
    if(-Not $setupExeFound)
    {
        $lstOldSqlVersionSetupExePaths = New-Object "System.Collections.Generic.List[string]"
    
        #SQL 2008
        $lstOldSqlVersionSetupExePaths.Add([System.IO.Path]::Combine($MSSQLpath, "100\Setup Bootstrap\Release\Setup.exe"))
    
        #SQL 2005
        $lstOldSqlVersionSetupExePaths.Add([System.IO.Path]::Combine($MSSQLpath, "90\Setup Bootstrap\Setup.exe"))
    
        #TODO: add strings to the array for even older versions of SQL (gulp).
    
        
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


function main(){
    if (-not $outDir.ToString().EndsWith("\")){
        $outDir = $($outDir.ToString()+'\' )
    }

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
        if ($outputAs -eq "file"){
            $result|ConvertTo-Json|Out-File -LiteralPath  $("$outDir"+'hostSqlFeatures.json')
        }else{
            $result|ConvertTo-Json
        }
    }else{
        Write-Host "No installed SQL Server features found." -ForegroundColor Yellow 
    }
}

main