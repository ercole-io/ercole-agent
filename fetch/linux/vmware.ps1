<#
        .SYNOPSIS
        Automated VMware cluster and VMs available in a list of provided vCenter servers.
        .DESCRIPTION
        Automated VMware cluster and VMs available in a list of provided vCenter servers.
        .EXAMPLE
        ./vmware.ps1 -s vms endpoint username password

        .EXAMPLE
        ./vmware.ps1 -s cluster endpoint username password
        .NOTES
        File Name  : cluster.ps1
        Author     : Riccardo Suardi - rsuardi@sorint.it
        Requires   : PowerShell Core, VMware PowerCLI
        .LINK
        https://sorint.it
        .Parameter s
        string, switch variable
        values accepted: cluster or vms
        ##
#>

param (
        [Parameter(Mandatory=$true)][string]$s,
        [Parameter(Mandatory=$true)][string]$endpoint,
        [Parameter(Mandatory=$true)][string]$username,
        [Parameter(Mandatory=$true)][string]$password
)

#Set-PowerCLIConfiguration -InvalidCertificateAction:Ignore
Connect-VIServer "$endpoint" -User "$username" -Password "$password" | Out-Null
New-VIProperty -Name NumCPU -ObjectType Cluster -Value {
                        $TotalPCPU = 0
                        $Args[0] | Get-VMHost | Foreach {
                                $TotalPCPU += $_.NumCPU
                        }
                        $TotalPCPU
        } `
        -Force -WarningAction:SilentlyContinue | Out-Null

New-VIProperty -Name NumSockets -ObjectType Cluster -Value {
                        $TotalPSOCKS = 0
                        $Args[0] | Get-VMHost | Foreach {
                                $TotalPSOCKS += $_.ExtensionData.Hardware.CpuInfo.NumCpuPackages
                        }
                        $TotalPSOCKS
        } `
        -Force -WarningAction:SilentlyContinue | Out-Null
switch ($s.ToUpper()) {
        "VMS" {
                # OUTPUT FORMAT: cluster name, vm name, guest os hostname
        	Get-VM | Select @{N="Cluster";E={if($c=Get-Cluster -VM $_){$c.name}else{Get-VMHost -VM $_}}}, Name, @{N="guestHostname";E={$_.ExtensionData.Guest.HostName}}, @{N="ESX Host";E={Get-VMHost -VM $_}} | ConvertTo-CSV | % { $_ -replace '"', ""}
	}
        "CLUSTER" {
                # OUTPUT FORMAT: cluster name, core sum, socket sum
                Get-Cluster | Select Name, NumCPU, NumSockets | ConvertTo-CSV | % { $_ -replace '"', ""}
		# OUTPUT (single node not in a Cluster): node name, core, socket
		Get-VMHost | ForEach-Object {
			$vmhost = $_
			$cluster = $vmhost | Get-Cluster
			if (-not $cluster) { # Check for hosts without clusters
				[PSCustomObject]@{
					Name = $vmhost.Name
					NumCPU = $_.NumCPU
					NumSockets = $_.ExtensionData.Hardware.CpuInfo.NumCpuPackages			
				}
			}
		} | Where-Object {$_.Name -ne $null} | ConvertTo-CSV | Select-Object -Skip 1 | % { $_ -replace '"', ""}
        }
        Default { Write-Host "wrong switch selection" }
}
Disconnect-VIServer $endpoint -Confirm:$false | Out-Null