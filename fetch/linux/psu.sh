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
DBV=$2
HOME=$3

if [ -z "$SID" ]; then
  >&2 echo "Missing SID parameter"
  exit 1
fi

if [ -z "$DBV" ]; then
  >&2 echo "Missing DBV parameter"
  exit 1
fi

if [ -z "$HOME" ]; then
  >&2 echo "Missing ORACLE_HOME parameter"
  exit 1
fi

USER=$4
PASSWORD=$5
if [ -z "$USER" ] && [ -z "$PASSWORD" ]; then
  SQLPLUS_CMD="sqlplus -S / as sysdba"
else
  SQLPLUS_CMD="sqlplus -S $USER/$PASSWORD"
fi

export ORAENV_ASK=NO 
export ORACLE_SID=$SID
export ORACLE_HOME=$HOME
export PATH=$HOME/bin:$PATH

LINUX_FETCHERS_DIR=$(dirname "$0")
FETCHERS_DIR="$(dirname "$LINUX_FETCHERS_DIR")"
ERCOLE_HOME="$(dirname "$FETCHERS_DIR")"

if [ $DBV -lt "12" ]; then
  $SQLPLUS_CMD @${ERCOLE_HOME}/sql/psu-1.sql
else
  $SQLPLUS_CMD @${ERCOLE_HOME}/sql/psu-2.sql
fi
