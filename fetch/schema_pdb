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
PDB=$3

if [ -z "$SID" ]; then
  echo "Missing SID parameter"
  exit 1
fi
if [ -z "$HOME" ]; then
  echo "Missing ORACLE_HOME parameter"
  exit 1
fi
if [ -z "$PDB" ]; then
  echo "Missing PDB parameter"
  exit 1
fi

ERCOLE_HOME=$(dirname "$0")
ERCOLE_HOME="$(dirname "$ERCOLE_HOME")"

export ORAENV_ASK=NO 
export ORACLE_SID=$SID
export ORACLE_HOME=$HOME
export PATH=$HOME/bin:$PATH

sqlplus -S "/ AS SYSDBA" @${ERCOLE_HOME}/sql/schema_pdb.sql $PDB