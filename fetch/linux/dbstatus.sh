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

SID=$1
HOME=$2

if [ -z "$SID" ]; then
  >&2 echo "Missing SID parameter"
  exit 1
fi
if [ -z "$HOME" ]; then
  >&2 echo "Missing ORACLE_HOME parameter"
  exit 1
fi

export ORAENV_ASK=NO 
export ORACLE_SID=$SID
export ORACLE_HOME=$HOME
export PATH=$HOME/bin:$PATH

RUNNING=$(ps -ef | grep pmon_$SID | grep -v grep | wc -l)
if [ "$RUNNING" -eq "0" ]; then
  echo "unreachable!"
  exit 0
fi

DB_STATUS=$(sqlplus -S / as sysdba << EOF
set pages 0 feedback off timing off
select open_mode from v\$database;
exit
EOF
)

echo -n $DB_STATUS
