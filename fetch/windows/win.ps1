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
	values accepted: host, fs, tab, dbversion, status, stats, db, mnt, tbs, schema, lic, patch, backup, addm and psu
		
	.Parameter d
	string, database name
	
	.Parameter v
	integer, database version

#>

param (
	[Parameter(Mandatory=$true)][string]$s,
	[Parameter()][string]$d,
	[Parameter()][int]$v,
	[Parameter()][string]$t,
	[Parameter()][string]$oraclepath
)


$pathAgent = Get-Location
cd $pathAgent\..\..

$constant = .50

#### DO NOT EDIT BELOW THIS LINE ####

$hname = $env:computername
$nfo_cpu	= gwmi win32_processor
$nfo_opsys 	= gwmi win32_operatingsystem
$nfo_sys	= gwmi win32_computersystem

function isVirtual {
	$ctr = 0
	$tiers = @('vmware','ovm','xen','virtual','hyper-v','citrix','innotek','oVirt')
	foreach ($tier in $tiers) {
		if ($nfo_sys.manufacturer -match $tier -or $nfo_sys.model -match $tier) { $ctr++ }
	}
	if($ctr) { return "Y" } else { return "N" }
}

function getType {
	$tiers = @('vmware','ovm','xen','virtual','hyper-v','citrix','innotek','oVirt')
	foreach ($tier in $tiers) {
		if ($nfo_sys.manufacturer -match $tier -or $nfo_sys.model -match $tier) { return $tier.ToUpper() }
	}
	return "PH"
}

function checkCommand($cmdname) {
    return [bool](Get-Command -Name $cmdname -ErrorAction SilentlyContinue)
}

function getSysinfo {
	$stmem	 = [math]::round($nfo_sys.totalphysicalmemory/1GB,0)
	$stvmem	 = [math]::round($nfo_opsys.totalvirtualmemorysize/1MB,0)
	Write-Host "Hostname:"$hname																	#hostname
	$crs  = 0
	$nfo_cpu | foreach { $_.numberofcores } | Foreach { $crs += $_}
	$sck = 0
	$nfo_cpu | foreach { $_.numberoflogicalprocessors } | Foreach { $sck += $_}
	if ($nfo_cpu -is [array]) {   #cpu model, cores and threads
		Write-Host "CPUModel:"$($nfo_cpu[0].name)  #is multisocket (pick the 1st label)
		Write-Host "CPUFrequency:"$($nfo_cpu[0].MaxClockSpeed)" MHz"
	} else {	#is singlesocket
		Write-Host "CPUModel:"$($nfo_cpu.name)
		Write-Host "CPUFrequency:"$($nfo_cpu.MaxClockSpeed)"MHz"
	}
	
	Write-Host "CPUCores:"$crs					
	Write-Host "CPUThreads:"$sck
	if ($(isVirtual) -eq "Y") {
		Write-Host "HardwareAbstraction: VIRT"
		Write-Host "CPUSockets:0" 						#sockets no
	} else {
		Write-Host "HardwareAbstraction: PH"
		Write-Host "CPUSockets:"$($nfo_cpu | foreach { $_.socketdesignation }).count 						#sockets no
	}
	Write-Host "ThreadsPerCore: 2"
	Write-Host "CoresPerSocket:"$($constant * $crs)
	Write-Host "HardwareAbstractionTechnology:"$(getType)
	Write-Host "Kernel: NT"														#kernel
	Write-Host "KernelVersion:"$($nfo_opsys.version)														#kernel version
	Write-Host "OS:"$($nfo_opsys.caption)		
	Write-Host "OSVersion:"$($nfo_opsys.caption)															#os#os
	Write-Host "MemoryTotal:"$stmem																	#total memory
	Write-Host "SwapTotal:"$stvmem																	#virtual memory
}

function getPartitions {
	$nfo_part = Get-WmiObject -Class win32_Volume | Where {$_.Drivetype -ne 5 -and $_.DriveLetter -ne $null} | select DeviceID, FileSystem, Capacity, FreeSpace, DriveLetter
	if (!$nfo_part) { Write-Warning "no partitions to list"; exit }
	foreach ($part in $nfo_part) {
		if ($part.DriveLetter) { write-host ("{0} {1} {2} {3} {4} {5}% {6}" -f $part.DeviceID,$part.FileSystem,$($part.Capacity / 1024),$($part.Capacity  / 1024-$part.FreeSpace / 1024),$($part.FreeSpace / 1024),$([math]::round((($part.Capacity / 1024-$part.FreeSpace / 1024)*100)/$part.Capacity / 1024,2)),$($part.DriveLetter)) }
	}
}

function getTab {
	$dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" }
	if (!$dbs) { Write "" } #no instances
	else {
		if ($dbs -is [array]) {
			foreach ($db in $dbs) {
				Write-Host $db.PathName.Split()[1]":"
			}
		}
		else {
			Write-Host $dbs.PathName.Split()[1]":"
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
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\stats.sql)) { Write-Warning "file stats.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @.\sql\stats.sql'
			if( $u -and $p ) { $ar = "-silent $u/$p @.\sql\stats.sql" }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getStatus {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\dbstatus.sql)) { Write-Warning "file dbstatus.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @.\sql\dbstatus.sql'
			if ($u -and $p) { $ar = "-silent $u/$p @.\sql\dbstatus.sql" }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDb {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$awr,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)

	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } }
	else { Write-Warning "missing arguments"; throw }

	if (!$dbs) { Write "" } #wrong or no instance
	else {
		if($awr -le 0){Write-Host "awr must be greater than zero"; throw}

		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\db.sql)) { Write-Warning "file db*.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\db.sql" '+ $awr
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\db.sql" '+ $awr }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbMount {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$true)][int]$v,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\dbmounted.sql)) { Write-Warning "file dbmounted.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\dbmounted.sql" '+$d+' '+$v
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\dbmounted.sql" '+$d+' '+$v }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbSchema {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\schema.sql)) { Write-Warning "file schema.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\schema.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\schema.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbTbs {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\ts.sql)) { Write-Warning "file ts.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\ts.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\ts.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbPartitionings {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\partitioning.sql)) { Write-Warning "file partitioning.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\partitioning.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\partitioning.sql" '+$d }
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
		[Parameter(Mandatory=$true)]$t,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID = $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\edition.sql)) { Write-Warning "file edition.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\edition.sql"'
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\edition.sql"' }
			Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow -RedirectStandardOutput $env:temp\datafile.xyz
			$edt = gc $env:temp\datafile.xyz
			rm $env:temp\datafile.xyz
		}
		$db_one = 'xOne'
		if (!((Test-Path .\sql\license.sql) -and (Test-Path .\sql\license-10.sql) -and (Test-Path .\sql\license_pluggable.sql))) { Write-Warning "file license*.sql unavailable!"; throw }
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
				10 { 
					$ar = '-silent / as sysdba @".\sql\license-10.sql" '+$crs+' '+$factor 
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\license-10.sql" '+$crs+' '+$factor }
					}
				11 { 
					$ar = '-silent / as sysdba @".\sql\license.sql" '+$crs+' '+$factor+' '+$db_one
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\license.sql" '+$crs+' '+$factor+' '+$db_one }
					}
				Default {
					if (!(Test-Path .\sql\pluggable.sql)) { Write-Warning "file pluggable.sql unavailable!"; throw }
					else {
						$arPDB = '-silent / as sysdba @".\sql\pluggable.sql"'
						if ($u -and $p) { $arPDB = '-silent '+"$u/$p"+' @".\sql\pluggable.sql"' }
						Start-Process $ohome\sqlplus -ArgumentList $arPDB -Wait -NoNewWindow -RedirectStandardOutput $env:temp\datafilePDB.xyz
						$idPDB = gc $env:temp\datafilePDB.xyz
						rm $env:temp\datafilePDB.xyz
						if ($idPDB == "TRUE") {
							$ar = '-silent / as sysdba @".\sql\license_pluggable.sql" '+$crs+' '+$factor+' '+$db_one
							if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\license_pluggable.sql" '+$crs+' '+$factor+' '+$db_one }
						} else {
							$ar = '-silent / as sysdba @".\sql\license.sql" '+$crs+' '+$factor+' '+$db_one
							if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\license.sql" '+$crs+' '+$factor+' '+$db_one }
						}
					}
				}
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
		[Parameter(Mandatory=$true)][int]$v,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\patch.sql) -and (Test-Path .\sql\patch-12.sql))) { Write-Warning "file patch*.sql unavailable!"; throw }
		else {
			switch ($v) {
				12 { 
					$ar = '-silent / as sysdba @".\sql\patch-12.sql" '+$d 
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\patch-12.sql" '+$d }
					}
				Default { 
					$ar = '-silent / as sysdba @".\sql\patch.sql" '+$d 
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\patch.sql" '+$d }
					}
				
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
		[Parameter(Mandatory=$true)][int]$v,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d -and $v) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\psu-1.sql) -and (Test-Path .\sql\psu-2.sql))) { Write-Warning "file psu.sql unavailable!"; throw }
		else {
			switch ($v) {
				10 { 
					$ar = '-silent / as sysdba @".\sql\psu-1.sql" '+$d 
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\psu-1.sql" '+$d }
					}
				11 { 
					$ar = '-silent / as sysdba @".\sql\psu-1.sql" '+$d
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\psu-1.sql" '+$d }
					}
				Default { 
					$ar = '-silent / as sysdba @".\sql\psu-2.sql" '+$d 
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\psu-2.sql" '+$d }
					}
			}			
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbBackup {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\backup_schedule.sql)) { Write-Warning "file backup_schedule.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\backup_schedule.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\backup_schedule.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbADDM {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\addm.sql)) { Write-Warning "file addm.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\addm.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\addm.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbSegmentAdvisor {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\segment_advisor.sql)) { Write-Warning "file segment_advisor.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\segment_advisor.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\segment_advisor.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getDbOpt {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!((Test-Path .\sql\opt.sql) -and (Test-Path .\sql\opt_pluggable.sql))) { Write-Warning "file opt.sql unavailable!"; throw }
		else {
			if (!(Test-Path .\sql\pluggable.sql)) { Write-Warning "file pluggable.sql unavailable!"; throw }
			else {
				$arPDB = '-silent / as sysdba @".\sql\pluggable.sql"'
				if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\pluggable.sql" '+$d }
				Start-Process $ohome\sqlplus -ArgumentList $arPDB -Wait -NoNewWindow -RedirectStandardOutput $env:temp\datafilePDB.xyz
				$idPDB = gc $env:temp\datafilePDB.xyz
				rm $env:temp\datafilePDB.xyz
				if ($idPDB == "TRUE") {
					$ar = '-silent / as sysdba @".\sql\opt_pluggable.sql" '+$d
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\opt_pluggable.sql" '+$d }
				} else {
					$ar = '-silent / as sysdba @".\sql\opt.sql" '+$d
					if ($u -and $p) { $ar = '-silent '+"$u/$p"+'  @".\sql\opt.sql" '+$d }
				}
				if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
					Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
				}
			}
		}
	}
}

function getServices {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "oracleservice" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\services.sql)) { Write-Warning "file services.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\services.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\services.sql" '+$d }
			if ( $dbs.state -eq "Running" -and $dbs.status -eq "OK" ) {
				Start-Process $ohome\sqlplus -ArgumentList $ar -Wait -NoNewWindow
			}
		}
	}
}

function getGrantsDba {
	param (
		[Parameter(Mandatory=$true)]$d,
		[Parameter(Mandatory=$false)]$u,
		[Parameter(Mandatory=$false)]$p
	)
	if ($d) { $dbs = gwmi -Class Win32_Service | ? { $_.name -match "grant_dba" -and $_.name -match $d } } else { Write-Warning "missing arguments"; throw }
	if (!$dbs) { Write "" } #wrong or no instance
	else {
		$ohome = ($dbs.PathName.Split()[0]).trim("ORACLE.EXE")
		$env:ORACLE_SID= $dbs.PathName.Split()[1]
		if (!(Test-Path .\sql\services.sql)) { Write-Warning "file services.sql unavailable!"; throw }
		else {
			$ar = '-silent / as sysdba @".\sql\grant_dba.sql" '+$d
			if ($u -and $p) { $ar = '-silent '+"$u/$p"+' @".\sql\grant_dba.sql" '+$d }
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
	"STATS"				{ getStats $d $u $p }
	"DBSTATUS"			{ getStatus $d $u $p }
	"DB"				{ getDb $d $v $u $p  }
	"DBMOUNTED"			{ getDbMount $d $v $u $p }
	"TABLESPACE"		{ getDbTbs $d $u $p }
	"SCHEMA"			{ getDbSchema $d $u $p }
	"LICENSE"			{ getDbLic $d $v $t $u $p }
	"PATCH"				{ getDbPatch $d $v $u $p }
	"PSU"				{ getDbPSU $d $v $u $p }
	"BACKUP"			{ getDbBackup $d $u $p }
	"ADDM"				{ getDbADDM $d $u $p }
	"SEGMENTADVISOR"	{ getDBSegmentAdvisor $d $u $p }
	"OPT"				{ getDbOpt $d $u $p }
	"SERVICES"			{ getServices $d $u $p }
	"GRANT_DBA"			{ getGrantsDba $d $u $p }
	"PARTITIONING"      { getDbPartitionings $d $u $p }
	Default				{ Write-Host "wrong switch selection" }
}
