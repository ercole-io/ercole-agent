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

# ASSUMPTIONS :
# MUST be run on physical server

#Check DBnodes and StorageServers
IFS=$'\n'
for i in $(ibhosts | sort -k 6); do
	host=$(echo $i | awk '{print $6}' | sed -e 's/"//' | sort -n)
	category=$(echo $i | awk '{print $7}')
	model="-"
	type="-"
	cpu="-"
	FCmode="-"
	cs="-"
	ping -c 2 ${host} >/dev/null
	if [ $? -ne 0 ]; then
		out="_SERVER_UNREACHABLE_"
	else
		ssh -o ConnectTimeout=2 -o BatchMode=yes -q ${host} exit
		if [ $? -ne 0 ]; then
			out="_CHECK_SSH_KEY_"
		else
			out=$(dcli -l root -c ${host} imageinfo -version | sed -e 's/://' | awk '{print $2}')
			case ${category} in
			*"S"*)
				type="DBServer"
				info=$(dcli -l root -c ${host} "dbmcli -e list dbserver detail" | egrep "Model|cpuCount|fanCount|fanStatus|memoryGB|powerCount|powerStatus|status|temperatureReading|temperatureStatus|cellsrvStatus|msStatus|rsStatus" | sed "s/${host}: //g")
				model=$(echo ${info} | awk '{print $12}')
				cpu=$(echo ${info} | awk '{print $2}')
				fanc=$(echo ${info} | awk '{print $4}')
				fans=$(echo ${info} | awk '{print $6}')
				if [[ $(ssh ${host} "ps -ef|grep xend|grep -v grep|wc -l") -gt 0 ]]; then
					mem=$(ssh ${host} "xm info|grep total_memory" | awk '{print $NF /1024}' | bc | awk -F"." '{print $1}')
				else
					mem=$(ssh ${host} "free -h |grep Mem" | awk '{print $2}' | sed -e 's/G//')
				fi
				#mem=$(ssh ${host} "if [ $(ps -ef|grep xend|grep -v grep|wc -l) -gt 0 ];then echo $(xm info|grep total_memory|awk '{print $NF}')/1024|bc ;else free -h |grep Mem|awk '{print $2}'|sed -e 's/G//';fi"|awk '{print $NF}')
				powerc=$(echo ${info} | awk '{print $14}')
				powers=$(echo ${info} | awk '{print $16}')
				stat=$(echo ${info} | awk '{print $18}')
				tempr=$(echo ${info} | awk '{print $20}')
				temps=$(echo ${info} | awk '{print $22}')
				ms=$(echo ${info} | awk '{print $24}')
				rs=$(echo ${info} | awk '{print $26}')
				;;
			*"C"*)
				type="StorageServer"
				info=$(dcli -l root -c ${host} "cellcli -e list cell detail" | egrep "Model|cpuCount|fanCount|fanStatus|flashCacheMode|memoryGB|powerCount|powerStatus|status|temperatureReading|temperatureStatus|cellsrvStatus|msStatus|rsStatus" | sed "s/${host}: //g")
				model=$(echo ${info} | awk '{print $14"_"$15"_"$16}')
				cpu=$(echo ${info} | awk '{print $2}')
				fanc=$(echo ${info} | awk '{print $4}')
				fans=$(echo ${info} | awk '{print $6}')
				mem=$(echo ${info} | awk '{print $18}')
				powerc=$(echo ${info} | awk '{print $20}')
				powers=$(echo ${info} | awk '{print $22}')
				stat=$(echo ${info} | awk '{print $24}')
				tempr=$(echo ${info} | awk '{print $26}')
				temps=$(echo ${info} | awk '{print $28}')
				cs=$(echo ${info} | awk '{print $30}')
				ms=$(echo ${info} | awk '{print $32}')
				rs=$(echo ${info} | awk '{print $34}')
				FCmode=$(echo ${info} | awk '{print $8}')
				;;

			esac
		fi
	fi
	echo "$host|||$type|||$model|||$out|||$cpu|||${mem}|||$stat|||$powerc|||$powers|||$fanc|||$fans|||$tempr|||$temps|||$cs|||$ms|||$rs|||$FCmode"
done

#Check InfiniBand Switches
for i in $(ibnodes | grep Switch | awk '{print $10}' | sort -n); do
	host=$i
	type="IBSwitch"
	model="-"
	cpu="-"
	mem="-"
	stat="-"
	powerc="-"
	powers="-"
	fanc="-"
	fans="-"
	tempr="-"
	temps="-"
	cs="-"
	ms="-"
	rs="-"
	FCmode="-"
	ping -c 2 ${host} >/dev/null
	if [ $? -ne 0 ]; then
		out="_SERVER_UNREACHABLE_"
	else
		ssh -o ConnectTimeout=2 -o BatchMode=yes -q ${host} exit
		if [ $? -ne 0 ]; then
			out="_CHECK_SSH_KEY_"
		else
			#info=$(dcli -l root -c ${host} 'version')
			out=$(dcli -l root -c ${host} 'version' | grep version | grep -v BIOS | awk -F":" '{print $NF}' | sed -e 's/ //')
			build=$(date -d"$(dcli -l root -c ${host} version | grep Build | awk '{print $4, $5, $6}')" +%y%m%d)
			model=$(dcli -l root -c ${host} 'version' | grep version | grep -v BIOS | awk '{print $2"_"$3"_"$4}')
			out=${out}.${build}
		fi
	fi
	echo "$host|||$type|||$model|||$out|||$cpu|||$mem|||$stat|||$powerc|||$powers|||$fanc|||$fans|||$tempr|||$temps|||$cs|||$ms|||$rs|||$FCmode"
done

#Example Output
#HOSTNAME|||SERVER_TYPE|||MODEL|||EXA_SW_VERSION|||CPU_ENABLED|||MEMORY|||STATUS|||POWER_COUNT|||POWER_STATUS|||FAN_COUNT|||FAN_STATUS|||TEMP_ACTUAL|||TEMP_STATUS|||CELLSRV_SERVICE|||MS_SERVICE|||RS_SERVICE|||FLASHCACHE_MODE
#-------------------------------------------------------------------------------------------
#fcax1hf1|||DBServer|||X8-2|||19.2.7.0.0.191012|||4/96|||766|||online|||2/2|||normal|||16/16|||normal|||22.0|||normal|||-|||running|||running|||-
#fcax1hf2|||DBServer|||X8-2|||19.2.7.0.0.191012|||4/96|||766|||online|||2/2|||normal|||16/16|||normal|||23.0|||normal|||-|||running|||running|||-
#fcax1sf1|||StorageServer|||X8-2L_High_Capacity|||19.2.7.0.0.191012|||64/64|||188|||online|||2/2|||normal|||8/8|||normal|||20.0|||normal|||running|||running|||running|||WriteBack
#fcax1sf2|||StorageServer|||X8-2L_High_Capacity|||19.2.7.0.0.191012|||64/64|||188|||online|||2/2|||normal|||8/8|||normal|||22.0|||normal|||running|||running|||running|||WriteBack
#fcax1sf3|||StorageServer|||X8-2L_High_Capacity|||19.2.7.0.0.191012|||64/64|||188|||online|||2/2|||normal|||8/8|||normal|||23.0|||normal|||running|||running|||running|||WriteBack
#fcax1ba0|||IBSwitch|||SUN_DCS_36p|||2.2.13-2.190326|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-
#fcax1bb0|||IBSwitch|||SUN_DCS_36p|||2.2.13-2.190326|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-|||-

#SOME INFO
#powerStatus: Status of the power. The value can be normal, warning, or critical.
#fanStatus: Status of the fan. The value can be normal, warning, or critical.
#temperatureStatus: Status of the temperature. The value can be normal, warning, or critical.
