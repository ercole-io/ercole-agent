#!/bin/sh

# Copyright (c) 2023 Sorint.lab S.p.A.
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

ORATAB_FILE=$1

ps -ef | grep pmon | grep -v grep| grep -v ASM| grep -v APX| grep -v FOMUT0_U| grep -v FOCSB0_T > /tmp/listPmonErcole.log

while read i
do PROC=$(echo $i| awk '{print $2}')
SID=$(echo $i | awk -F "_" '{print $NF}')
ENTRY=$(echo "$SID:$(pwdx $PROC| awk -F : '{print $2}')"| tr -d ' ')
if [[ $(grep $SID "$ORATAB_FILE"|wc -l) -lt 1 ]] || ([[ $(grep "#"$SID "$ORATAB_FILE"|wc -l) -eq 1 ]] && [[ $(grep $SID "$ORATAB_FILE"|wc -l) -lt 2 ]]); then
	echo ${ENTRY%????}:N
fi
done </tmp/listPmonErcole.log