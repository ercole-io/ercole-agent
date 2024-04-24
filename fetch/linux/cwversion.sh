#!/bin/bash

# Copyright (c) 2024 Sorint.lab S.p.A.
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

if [[ -f /etc/oraInst.loc ]]; then
  INVENTORY=$(grep inventory_loc /etc/oraInst.loc|awk -F['='] '{print $2}')
  GI_HOME=$(grep 'CRS="true"' ${INVENTORY}/ContentsXML/inventory.xml|awk -F['='] '{print $3'}|awk -F['"'] '{print $2}')
  if [[ -d $GI_HOME ]]; then
    CHECK_ASM=$(ps -ef|grep asm_pmon|grep -iv grep|wc -l)
    if [[ $CHECK_ASM == 1 ]]; then
      CWVERSION=$(${GI_HOME}/bin/srvctl -version| awk '{print $3}')
	  CWVERSION=$(echo $CWVERSION | grep -oE '[0-9]+[.]+[0-9]+' | head -1)
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

echo -n "CWVersion: $CWVERSION
"