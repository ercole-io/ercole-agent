#!/bin/sh

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
PROC=$1
RUNNING_INSTANCE=$2


SID=$(echo $RUNNING_INSTANCE | awk -F "ora_pmon_" '{print $NF}')
ENTRY=$(echo "$SID:$(sudo pwdx $PROC| awk -F : '{print $2}')"| tr -d ' ')
echo ${ENTRY%????}:N