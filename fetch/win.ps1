<#

	.SYNOPSIS
	Automated system informations and Oracle database informations extraction.

	.DESCRIPTION
	Automated system informations and Oracle database informations extraction for Windows environments called from the server side.
	Output is redirected on STDOUT in HTML format.

	.EXAMPLE
	./win.ps1 -s lic -d ORATEST -v 12
	
	.EXAMPLE
	./win.ps1 -s host

	.NOTES
	File Name  : win.ps1  
	Author     : Riccardo Suardi - rsuardi@sorint.it 
	Requires   : PowerShell V4 

	.LINK
	https://ercole.io

	.Parameter s
	string, switch variable
	values accepted: host, fs, tab, dbversion, status, stats, db, mnt, feature, tbs, schema, lic, patch, backup, addm and psu
		
	.Parameter d
	string, database name
	
	.Parameter v
	integer, database version

#>

param (
	[Parameter(Mandatory=$true)][string]$s,
	[Parameter()][string]$d,
	[Parameter()][int]$v,
    [Parameter()][string]$t
)

$constant = .50

#### DO NOT EDIT BELOW THIS LINE ####

$hname = $env:computername
$nfo_cpu	= gwmi win32_processor
$nfo_opsys 	= gwmi win32_operatingsystem
$nfo_sys	= gwmi win32_computersystem

function isVirtual {
	$ctr = 0
	$tiers = @('vmware','ovm','xen','virtual','hyper-v','citrix','innotek')
	foreach ($tier in $tiers) {
		if ($nfo_sys.manufacturer -match $tier -or $nfo_sys.model -match $tier) { $ctr++ }
	}
	if($ctr) { return "Y" } else { return "N" }
}

function checkCommand($cmdname)
{
    return [bool](Get-Command -Name $cmdname -ErrorAction SilentlyContinue)
}
function getSysinfo {
	$stmem	 = [math]::round($nfo_sys.totalphysicalmemory/1GB,0)
	$stvmem	 = [math]::round($nfo_opsys.totalvirtualmemorysize/1MB,0)
	Write-Host "Hostname:"$hname																	#hostname
	$crs  = 0
	$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
	$sck = 0
	$sck += $($nfo_cpu | foreach { $_.numberoflogicalprocessors } | Foreach { $sck += $_})
	if ($nfo_cpu -is [array]) {   #cpu model, cores and threads
		Write-Host "CPUModel:"$($nfo_cpu[0].name)  #is multisocket (pick the 1st label)
	} else {	#is singlesocket
		Write-Host "CPUModel:"$($nfo_cpu.name)
	}
	Write-Host "CPUCores:"$crs					
	Write-Host "CPUThreads:"$sck
	if ($(isVirtual) -eq "Y") {
		Write-Host "Socket:0" 						#sockets no
	} else {
		Write-Host "Socket:"$($nfo_cpu | foreach { $_.socketdesignation }).count 						#sockets no
	}
	Write-Host "Virtual:"$(isVirtual)																#virtual/physical
	Write-Host "Kernel:"$($nfo_opsys.version)														#kernel version
	Write-Host "OS:"$($nfo_opsys.caption)															#os
	Write-Host "MemoryTotal:"$stmem																	#total memory
	Write-Host "SwapTotal:"$stvmem																	#virtual memory
	Write-Host "OracleCluster: N"
	Write-Host "VeritasCluster: N"
	Write-Host "SunCluster: N"
	Write-Host "AixCluster: N"
}

function getPartitions {
	$nfo_part = Get-WmiObject -Class win32_Volume | Where {$_.Drivetype -ne 5 -and $_.DriveLetter -ne $null} | select DeviceID, FileSystem, Capacity, FreeSpace, DriveLetter
	if (!$nfo_part) { Write-Warning "no partitions to list"; exit }
	foreach ($part in $nfo_part) {
		if ($part.DriveLetter) { write-host ("{0} {1} {2} {3} {4} {5}% {6}" -f $part.DeviceID,$part.FileSystem,$($part.Capacity),$($part.Capacity-$part.FreeSpace),$($part.FreeSpace),$([math]::round((($part.Capacity-$part.FreeSpace)*100)/$part.Capacity,2)),$($part.DriveLetter)) }
	}
}

function getTab {
	$dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" }
	if (!$dbs) { Write "" } #no instances
	else {
		if ($dbs -is [array]) {
			foreach ($db in $dbs) {
				write $db.PathName.Split()[1]
			}
		}
		else {
			write $dbs.PathName.Split()[1]
		}
	}
}

function getVer {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else { write $(gci $dbs.PathName.Split()[0]).VersionInfo.fileversion.split(" ")[0] }
}

function getStats {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\stats.sql)) { Write-Warning "file stats.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @.\sql\stats.sql'
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getStatus {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\dbstatus.sql)) { Write-Warning "file dbstatus.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @.\sql\dbstatus.sql'
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDb {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\db.sql)) { Write-Warning "file db*.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\db.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbMount {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\dbmounted.sql)) { Write-Warning "file dbmounted.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\dbmounted.sql" '+$d+' '+$v
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbFeature {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\feature.sql) -and (Test-Path .\sql\feature-10.sql))) { Write-Warning "file feature*.sql unavailable!"; throw }
		else {
			switch ($v) {
				10 { $ar = '-silent / as sysdba @".\sql\feature-10.sql" '+$d }
				Default { $ar = '-silent / as sysdba @".\sql\feature.sql" '+$d }
			}
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbSchema {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\schema.sql)) { Write-Warning "file schema.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\schema.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbTbs {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\ts.sql)) { Write-Warning "file ts.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\ts.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbLic {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v,
		[Parameter(Mandatory=$true)]$t
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID = $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\edition.sql)) { Write-Warning "file edition.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\edition.sql"'
			Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow -RedirectStandardOutput $env:temp\datafile.xyz
			$edt = gc $env:temp\datafile.xyz
			rm $env:temp\datafile.xyz
		}
		#if (!(Test-Path .\sql\dbone.sql)) { Write-Warning "file dbone.sql unavailable!"; throw }
		#else {
		#	$ar = '-silent / as sysdba @".\sql\dbone.sql"'
		#	Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow -RedirectStandardOutput $env:temp\datafile.xyz
		#	$db_one = gc $env:temp\datafile.xyz
		#	$db_one = 'x' + $db_one
		#	rm $env:temp\datafile.xyz
		#}
		$db_one = 'xOne'
		if (!((Test-Path .\sql\license.sql) -and (Test-Path .\sql\license-10.sql))) { Write-Warning "file license*.sql unavailable!"; throw }
		else {
			switch (isVirtual) {
				"Y" {
					switch ($edt.toUpper()) {
						"ENT" { 
							$crs = 0
							$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
							$core_factor = $constant * $crs
							$factor = $constant
						}
						"EXE" { 
							$crs = 0
							$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
							$core_factor = $constant * $crs
							$factor = $constant
						}
						"STD" { 
							$crs = 0
							$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
							$core_factor = $constant * $crs
							$factor = $constant
						}
					}
				}
				"N" {
					switch ($edt.toUpper()) {
						"ENT" {
							$crs = 0
							$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
							$core_factor = $constant * $crs
							$factor = $constant
						}
						"EXE" { 
							$crs = 0
							$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
							$core_factor = $constant * $crs
							$factor = $constant
						}
						"STD" {
							$crs = $($nfo_cpu | foreach { $_.socketdesignation }).count
							$core_factor = $constant * $crs
							$factor = $constant * $crs
						}
					}
				}
			}
			switch ($edt.toUpper()) {
				"EXE" { 
					Write-Host "Oracle EXE;`t"$core_factor";" 
					Write-Host "Oracle ENT;;" 
					Write-Host "Oracle STD;;" 
				}
				"ENT" {
					Write-Host "Oracle EXE;;" 
					Write-Host "Oracle ENT;`t"$core_factor";" 
					Write-Host "Oracle STD;;"  
				}
				"STD" {
					Write-Host "Oracle EXE;;" 
					Write-Host "Oracle ENT;;" 
					Write-Host "Oracle STD;`t"$core_factor";" 
				}
			}
			switch ($v) {
				10 { $ar = '-silent / as sysdba @".\sql\license-10.sql" '+$crs+' '+$factor }
				Default { $ar = '-silent / as sysdba @".\sql\license.sql" '+$crs+' '+$factor+' '+$db_one}
			}
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbPatch {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\patch.sql) -and (Test-Path .\sql\patch-12.sql))) { Write-Warning "file patch*.sql unavailable!"; throw }
		else {
			switch ($v) {
				12 { $ar = '-silent / as sysdba @".\sql\patch-12.sql" '+$d }
				Default { $ar = '-silent / as sysdba @".\sql\patch.sql" '+$d }
				
			}
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbPSU {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\psu-1.sql) -and (Test-Path .\sql\psu-2.sql))) { Write-Warning "file psu.sql unavailable!"; throw }
		else {
			switch ($v) {
				10 { $ar = '-silent / as sysdba @".\sql\psu-1.sql" '+$d }
				11 { $ar = '-silent / as sysdba @".\sql\psu-1.sql" '+$d }
				Default { $ar = '-silent / as sysdba @".\sql\psu-2.sql" '+$d }
			}			
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbBackup {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\backup_schedule.sql)) { Write-Warning "file backup_schedule.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\backup_schedule.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbADDM {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\addm.sql)) { Write-Warning "file addm.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\addm.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbSegmentAdvisor {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\segment_advisor.sql)) { Write-Warning "file segment_advisor.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\segment_advisor.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbOpt {
	param (
		[Parameter(Mandatory=$true)]$d
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\opt.sql)) { Write-Warning "file opt.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\opt.sql" '+$d
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

switch($s.ToUpper()) {
	"HOST"				{ getSysinfo }
	"FILESYSTEM"		{ getPartitions }
	"ORATAB"			{ getTab }
	"DBVERSION"			{ getVer $d }
	"STATS"				{ getStats $d }
	"DBSTATUS"			{ getStatus $d }
	"DB"				{ getDb $d $v }
	"DBMOUNTED"			{ getDbMount $d $v }
	"FEATURE"			{ getDbFeature $d $v }
	"TABLESPACE"		{ getDbTbs $d }
	"SCHEMA"			{ getDbSchema $d }
	"LICENSE"			{ getDbLic $d $v $t }
	"PATCH"				{ getDbPatch $d $v }
	"PSU"				{ getDbPSU $d $v }
	"BACKUP"			{ getDbBackup $d }
	"ADDM"				{ getDbADDM $d }
	"SEGMENTADVISOR"	{ getDBSegmentAdvisor $d }
	"OPT"				{ getDbOpt $d }
	Default				{ Write-Host "wrong switch selection" }
}
