#!/bin/bash

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

if [[ $(echo $CPU_MODEL | grep -c "@") -gt 0 ]]; then
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
if [[ $CORES_PER_SOCKET -le 0 ]]; then
  CORES_PER_SOCKET=1
fi

KERNEL=$(uname --kernel-name)
KERNEL_VERSION=$(uname --kernel-release)

if [[ -f /etc/redhat-release ]]; then
  OS="Red Hat Enterprise Linux"
else
  OS=$(grep "^NAME=" /etc/os-release | awk -F\= '{gsub(/"/,"",$2);print $2}')
fi

if [[ -z $OS ]]; then
  OS="unknown"
fi

if [[ -f /etc/os-release ]]; then
  OS_VERSION=$(grep "^VERSION_ID=" /etc/os-release | awk -F\= '{gsub(/"/,"",$2);print $2}')
elif [[ -f /etc/redhat-release ]]; then
    OS_VERSION=$(cat /etc/redhat-release | rev | cut -d' ' -f2 | rev)
fi

if [[ -z $OS_VERSION ]]; then
  OS_VERSION="unknown"
fi


MEM_TOTAL=$(echo "$(($(free -k | grep Mem | awk -F ' ' '{print $2}') / 1024 / 1024))")
SWP_TOTAL=$(echo "$(($(free -k | grep Swap | awk -F ' ' '{print $2}') / 1024 / 1024))")

CHECK_TYPE_SERVER_OVM_DMESG=$(dmesg | grep OVM | wc -l)
CHECK_TYPE_SERVER_OVM_LOG=$(grep OVM /var/log/dmesg | wc -l)
CHECK_TYPE_SERVER_KVM_DMESG=$(dmesg | grep -i "Hypervisor detected: KVM" | wc -l)
CHECK_TYPE_SERVER_KVM_LOG=$(grep -i "Hypervisor detected: KVM" /var/log/dmesg | wc -l)
CHECK_TYPE_SERVER_VMWARE=$(dmesg | grep VMware | wc -l)
CHECK_TYPE_SERVER_VMWARE_LOG=$(grep VMware /var/log/dmesg* | wc -l)
CHECK_TYPE_SERVER_HYPERV=$(dmesg | grep HyperV | wc -l)
CHECK_TYPE_SERVER_HYPERV_LOG=$(grep HyperV /var/log/dmesg* | wc -l)
CHECK_TYPE_SERVER_HPUX=0
CHECK_TYPE_SERVER_HYPERVISOR=$(grep ^flags /proc/cpuinfo | grep hypervisor | wc -l)

if [ "$CHECK_TYPE_SERVER_OVM_DMESG" -gt 0 ] || [ "$CHECK_TYPE_SERVER_OVM_LOG" -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=OVM
  VIRTUAL=Y
elif [ $CHECK_TYPE_SERVER_KVM_DMESG -gt 0 ] || [ "$CHECK_TYPE_SERVER_KVM_LOG" -gt 0 ]; then
  HARDWARE_ABSTRACTION_TECHNOLOGY=KVM
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

if [[ -f /etc/oraInst.loc ]]; then
  INVENTORY=$(grep inventory_loc /etc/oraInst.loc|awk -F['='] '{print $2}')
  GI_HOME=$(grep 'CRS="true"' ${INVENTORY}/ContentsXML/inventory.xml|awk -F['='] '{print $3'}|awk -F['"'] '{print $2}')
  if [[ -d $GI_HOME ]]; then
    CHECK_ASM=$(ps -ef|grep asm_pmon|grep -iv grep|wc -l)
    if [[ $CHECK_ASM == 1 ]]; then
      CWVERSION=$(${GI_HOME}/bin/asmcmd showversion| awk '{print $4}')
    else
      CHECK_PSU=$(${GI_HOME}/OPatch/opatch lspatches|grep 'Database Patch Set Update'|awk -F[:] '{print $2}'|awk -F ['('] {'print $1'})
      CHECK_RU=$(${GI_HOME}/OPatch/opatch lspatches|grep 'Database Release Update'|awk -F[:] '{print $2}'|awk -F ['('] {'print $1'})
      if [ ! -z $CHECK_PSU ] && [ -z $CHECK_RU ]; then 
        CWVERSION=$CHECK_PSU
      elif [ -z $CHECK_PSU ] && [ ! -z $CHECK_RU ]; then
        CWVERSION=$CHECK_RU
      else
        CWVERSION=" "
      fi
    fi
  else
    CWVERSION=" "
  fi
else
  CWVERSION=" "
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
CWVersion: $CWVERSION
"