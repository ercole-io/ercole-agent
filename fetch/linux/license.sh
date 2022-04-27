#!/bin/sh

# Copyright (c) 2022 Sorint.lab S.p.A.
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
TYPE=$3
HOME=$4
FACTOR=$5

if [ -z "$SID" ]; then
    echo >&2 "Missing SID parameter"
    exit 1
fi

if [ -z "$TYPE" ]; then
    echo >&2 "Missing type parameter"
    exit 1
fi

if [ -z "$DBV" ]; then
    echo >&2 "Missing DBV parameter"
    exit 1
fi

if [ -z "$HOME" ]; then
    echo >&2 "Missing ORACLE_HOME parameter"
    exit 1
fi

if [ -z "$FACTOR" ]; then
    echo >&2 "Missing FACTOR parameter"
    exit 1
fi

LINUX_FETCHERS_DIR=$(dirname "$0")
FETCHERS_DIR="$(dirname "$LINUX_FETCHERS_DIR")"
ERCOLE_HOME="$(dirname "$FETCHERS_DIR")"

export ORAENV_ASK=NO
export ORACLE_SID=$SID
export ORACLE_HOME=$HOME
export PATH=$HOME/bin:$PATH

DB_VERSION=$(
    sqlplus -S / as sysdba <<EOF
set pages 0 feedback off timing off
select (case when UPPER(banner) like '%EXTREME%' then 'EXE' when UPPER(banner) like '%ENTERPRISE%' then 'ENT' else 'STD' end) as versione from v\$version where rownum=1;
exit
EOF
)

DB_NAME=$(
    sqlplus -S / as sysdba <<EOF
set pages 0 feedback off timing off
select name from v\$database;
exit
EOF
)

DB_ONE=x$(
    sqlplus -S / as sysdba <<EOF
set pages 0 feedback off timing off
HOST srvctl config database -d $DB_NAME |grep -o One
exit
EOF
)

CPU_THREADS=$(grep processor /proc/cpuinfo | wc -l)
THREAD_FACTOR=$(echo "$FACTOR"/2 | bc --mathlib)

if [[ "$TYPE" == 'OVM' || "$TYPE" == 'VMWARE' || "$TYPE" == 'VMOTHER' || "$TYPE" == 'KVM' ]]; then
    if [[ $DB_VERSION == 'EXE' || $DB_VERSION == 'ENT' ]]; then
        LICENSES=$(echo "$THREAD_FACTOR"*$CPU_THREADS | bc)
    elif [[ $DB_VERSION == 'STD' ]]; then
        LICENSES=1
    fi
elif [ $TYPE == 'PH' ]; then
    if [[ $DB_VERSION == 'EXE' || $DB_VERSION == 'ENT' ]]; then
        LICENSES=$(echo "$THREAD_FACTOR"*$CPU_THREADS | bc)
    elif [[ $DB_VERSION == 'STD' ]]; then
        LICENSES=$(cat /proc/cpuinfo | grep -i "physical id" | sort -n | uniq | wc -l)
    fi
fi

if [[ $DB_VERSION == 'EXE' ]]; then
    echo "Oracle EXE; $LICENSES;"
else
    echo "Oracle EXE;;"
fi
if [[ $DB_VERSION == 'ENT' ]]; then
    echo "Oracle ENT; $LICENSES;"
else
    echo "Oracle ENT;;"
fi
if [[ $DB_VERSION == 'STD' ]]; then
    echo "Oracle STD; $LICENSES;"
else
    echo "Oracle STD;;"
fi

if [ $DBV == "10" ] || [ $DBV == "9" ]; then
    sqlplus -S "/ AS SYSDBA" @${ERCOLE_HOME}/sql/license-10.sql $CPU_THREADS "$THREAD_FACTOR"
elif [ $DBV == "11" ]; then 
    sqlplus -S "/ AS SYSDBA" @${ERCOLE_HOME}/sql/license.sql $CPU_THREADS "$THREAD_FACTOR" $DB_ONE
else
IS_PDB=$(
    sqlplus -S / as sysdba <<EOF
set pages 0 feedback off timing off
SELECT CASE WHEN COUNT(*) > 0 THEN 'TRUE' WHEN count(*) = 0 THEN 'FALSE' END FROM v\$pdbs;
exit
EOF
)

if [ $IS_PDB == "TRUE" ]; then
    sqlplus -S "/ AS SYSDBA" @${ERCOLE_HOME}/sql/license_pluggable.sql $CPU_THREADS "$THREAD_FACTOR" $DB_ONE
else
    sqlplus -S "/ AS SYSDBA" @${ERCOLE_HOME}/sql/license.sql $CPU_THREADS "$THREAD_FACTOR" $DB_ONE
fi
fi