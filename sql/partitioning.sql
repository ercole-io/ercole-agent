-- Copyright (c) 2022 Sorint.lab S.p.A.

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program.  If not, see <http://www.gnu.org/licenses/>.

set lines 8000 pages 0 feedback off verify off timing off
set colsep "|||" 

column OWNER for a20 
column SEGMENT_NAME for a40 
column PARTITION_NAME for a40 

select 
    owner,
    segment_name,
    count(*),
    sum(bytes/1024/1024) "MB"
from dba_segments 
where partition_name is not null 
and owner not in ('SYS','AUDSYS','SYSTEM','SYSBACKUP','SYSDG','SYSKM','OUTLN','GSMADMIN_INTERNAL', 'GSMUSER','DIP','XS$NULL','ORACLE_OCM','DBSNMP','APPQOSSYS','ANONYMOUS','XDB', 'GSMCATUSER','WMSYS','OJVMSYS','CTXSYS','ORDDATA','ORDSYS','ORDPLUGINS', 'SI_INFORMTN_SCHEMA','MDSYS','OLAPSYS','MDDATA','SPATIAL_WFS_ADMIN_USR', 'SPATIAL_CSW_ADMIN_USR','LBACSYS','APEX_040200','APEX_PUBLIC_USER','FLOWS_FILES', 'DVSYS','DVF','SCOTT','EXFSYS','XS\$NULL','CMDB_USR','CPFI0_APPPRD','SYSMAN','SYSADMIN', 'MGMT_VIEW','DMSYS','WCCREPUSER','WCCUSER','WFADMIN','OWF_MGR','KWALKER','BLEWIS','CDOUGLAS','SPIERSON')
group by segment_name,owner;

exit
