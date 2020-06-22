#!/bin/sh

# Copyright (c) 2019 Sorint.lab S.p.A.
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

HOSTNAME=$(hostname)

CPU_MODEL=$(grep "model name" /proc/cpuinfo | sort -u | awk -F ': ' '{print $2}')

if [[ $(echo $CPU_MODEL | grep -c "@") ]]; then
  CPU_FREQUENCY=$(echo $CPU_MODEL | awk -F '@' '{print $2}')
else
  CPU_FREQUENCY=$(grep "cpu MHz" /proc/cpuinfo | head -n 1 | awk -F ":" '{ print $NF }')Mhz
fi

CPU_MODEL=$(echo $CPU_MODEL | awk -F '@' '{print $1}')

CPU_THREADS=$(grep processor /proc/cpuinfo | wc -l)
if [[ $CPU_THREADS = '1' ]]; then
  CPU_CORES=1
else
  CPU_CORES=$(($CPU_THREADS / 2))
fi

THREADS_PER_CORE=$(expr $CPU_THREADS / $CPU_CORES)

if [[ $VIRTUAL = 'Y' ]]; then
  CPU_SOCKETS=1
else
  CPU_SOCKETS=$(cat /proc/cpuinfo | grep -i "physical id" | sort -n | uniq | wc -l)
fi

CORES_PER_SOCKET=$(expr $CPU_CORES / $CPU_SOCKETS)

KERNEL=$(uname --kernel-name)
KERNEL_VERSION=$(uname --kernel-release)

OS=$(cat /etc/redhat-release)
if [[ $? != 0 ]]; then
  OS=$(cat /etc/SuSE-release | head -1)
fi

OS=$(cat /etc/os-release | grep "^NAME=" | awk -F\= '{gsub(/"/,"",$2);print $2}')
OS_VERSION=$(cat /etc/os-release | grep "^VERSION_ID=" | awk -F\= '{gsub(/"/,"",$2);print $2}')

MEM_TOTAL=$(echo "$(($(free -k | grep Mem | awk -F ' ' '{print $2}') / 1024 / 1024))")
SWP_TOTAL=$(echo "$(($(free -k | grep Swap | awk -F ' ' '{print $2}') / 1024 / 1024))")

CHECK_SUN_CLUSTER=$(ps -eo cmd | grep -v grep | grep "/rpc.pmfd\b" | wc -l)
if [ $CHECK_SUN_CLUSTER = 1 ]; then
  SUN_CLUSTER=Y
else
  SUN_CLUSTER=N
fi

CHECK_VERITAS_CLUSTER=$(ps -eo cmd | grep -v grep | grep "/had\b" | wc -l)
if [ $CHECK_VERITAS_CLUSTER = 1 ]; then
  VERITAS_CLUSTER=Y
else
  VERITAS_CLUSTER=N
fi

CHECK_ORACLE_CLUSTER=$(ps -eo cmd | grep -v grep | grep "/crsd\b" | wc -l)
if [ $CHECK_ORACLE_CLUSTER -gt 0 ]; then
  ORACLE_CLUSTER=Y
else
  ORACLE_CLUSTER=N
fi

CHECK_TYPE_SERVER_OVM_DMESG=$(dmesg | grep OVM | wc -l)
CHECK_TYPE_SERVER_OVM_LOG=$(grep OVM /var/log/dmesg | wc -l)
CHECK_TYPE_SERVER_VMWARE=$(dmesg | grep VMware | wc -l)
CHECK_TYPE_SERVER_VMWARE_LOG=$(grep VMware /var/log/dmesg* | wc -l)
CHECK_TYPE_SERVER_HYPERV=$(dmesg | grep HyperV | wc -l)
CHECK_TYPE_SERVER_HYPERV_LOG=$(grep HyperV /var/log/dmesg* | wc -l)
CHECK_TYPE_SERVER_HPUX=0
CHECK_TYPE_SERVER_HYPERVISOR=$(grep ^flags /proc/cpuinfo | grep hypervisor | wc -l)

if [ "$CHECK_TYPE_SERVER_OVM_DMESG" -gt 0 ] || [ "$CHECK_TYPE_SERVER_OVM_LOG" -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=OVM
  VIRTUAL=Y
elif [ $CHECK_TYPE_SERVER_VMWARE -gt 0 ] || [ "$CHECK_TYPE_SERVER_VMWARE_LOG" -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=VMWARE
  VIRTUAL=Y
elif [ $CHECK_TYPE_SERVER_HYPERV -gt 0 ] || [ "$CHECK_TYPE_SERVER_HYPERV_LOG" -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=HYPERV
  VIRTUAL=Y
elif [ $CHECK_TYPE_SERVER_HYPERVISOR -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=VMOTHER
  VIRTUAL=Y
else
  HARDWARE_ABSTRACTION_TECHNOLOGY=PH
  VIRTUAL=N
fi

if [ $VIRTUAL == "Y" ]; then
  HARDWARE_ABSTRACTION="VIRT"
else
  HARDWARE_ABSTRACTION="PH"
fi

echo -n "Hostname: $HOSTNAME
CPUModel: $CPU_MODEL 
CPUFrequency: $CPU_FREQUENCY
CPUSockets: $CPU_SOCKETS
CPUCores: $CPU_CORES
CPUThreads: $CPU_THREADS
ThreadsPerCore: $THREADS_PER_CORE
CoresPerSocket: $CORES_PER_SOCKET 
HardwareAbstraction: $HARDWARE_ABSTRACTION
HardwareAbstractionTechnology: $HARDWARE_ABSTRACTION_TECHNOLOGY
Kernel: $KERNEL
KernelVersion: $KERNEL_VERSION
OS: $OS
OSVersion: $OS_VERSION
MemoryTotal: $MEM_TOTAL
SwapTotal: $SWP_TOTAL
OracleCluster: $ORACLE_CLUSTER
VeritasCluster: $VERITAS_CLUSTER
SunCluster: $SUN_CLUSTER
AixCluster: N"
