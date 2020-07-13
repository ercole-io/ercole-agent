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

CHECK_ORACLE_CLUSTERWARE=$(ps -eo cmd | grep -v grep | grep "/crsd\b" | wc -l)
if [ $CHECK_ORACLE_CLUSTERWARE -gt 0 ]; then
    ORACLE_CLUSTERWARE=Y
else
    ORACLE_CLUSTERWARE=N
fi

CHECK_VERITAS_CLUSTER_SERVER=$(ps -eo cmd | grep -v grep | grep "/had\b" | wc -l)
if [ $CHECK_VERITAS_CLUSTER_SERVER = 1 ]; then
    VERITAS_CLUSTER_SERVER=Y
else
    VERITAS_CLUSTER_SERVER=N
fi

CHECK_SUN_CLUSTER=$(ps -eo cmd | grep -v grep | grep "/rpc.pmfd\b" | wc -l)
if [ $CHECK_SUN_CLUSTER = 1 ]; then
    SUN_CLUSTER=Y
else
    SUN_CLUSTER=N
fi

echo -n "OracleClusterware: $ORACLE_CLUSTERWARE
VeritasClusterServer: $VERITAS_CLUSTER_SERVER
SunCluster: $SUN_CLUSTER
"
