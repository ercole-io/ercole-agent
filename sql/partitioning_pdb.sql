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

alter session set container=&1;

column OWNER for a20 
column SEGMENT_NAME for a40 
column PARTITION_NAME for a40 

with users_list as (
select username from dba_users where username not in ('SYS','AUDSYS','SYSTEM','SYSBACKUP','SYSDG','SYSKM','OUTLN','GSMADMIN_INTERNAL', 'GSMUSER','DIP','XS$NULL','ORACLE_OCM','DBSNMP','APPQOSSYS','ANONYMOUS','XDB', 'GSMCATUSER','WMSYS','OJVMSYS','CTXSYS','ORDDATA','ORDSYS','ORDPLUGINS', 'SI_INFORMTN_SCHEMA','MDSYS','OLAPSYS','MDDATA','SPATIAL_WFS_ADMIN_USR', 'SPATIAL_CSW_ADMIN_USR','LBACSYS','APEX_040200','APEX_PUBLIC_USER','FLOWS_FILES', 'DVSYS','DVF','SCOTT','EXFSYS','XS\$NULL','CMDB_USR','CPFI0_APPPRD','SYSMAN','SYSADMIN', 'MGMT_VIEW','DMSYS','WCCREPUSER','WCCUSER','WFADMIN','OWF_MGR','KWALKER','BLEWIS','CDOUGLAS','SPIERSON')
),
--Backing table useful in the minus part
segment_lobs as (
(select l.owner,l.segment_name as segment_name from dba_lobs l inner join
dba_segments s on l.owner=s.owner and l.table_name=s.segment_name 
where s.segment_type in ('TABLE PARTITION') and s.owner in (select username from users_list)
group by l.owner,l.segment_name
union
select l.owner,l.index_name as segment_name from dba_lobs l inner join
dba_segments s on l.owner=s.owner and l.table_name=s.segment_name 
where s.segment_type in ('TABLE PARTITION') and s.owner in (select username from users_list)
group by l.owner,l.index_name)
)
--Retrieves partitioned segments
select 
    owner,
    segment_name,
    count(*),
    sum(bytes/1024/1024) "MB"
from dba_segments 
where partition_name is not null 
and owner in (select username from users_list)
group by segment_name,owner
minus
--Removes partitioned segments created automatically by using lob columns (not interesting)
select 
    dba_segments.owner,
    dba_segments.segment_name,
    count(*),
    sum(bytes/1024/1024) "MB"
from dba_segments inner join segment_lobs 
on dba_segments.owner=segment_lobs.owner and dba_segments.segment_name=segment_lobs.segment_name
where dba_segments.partition_name is not null 
and dba_segments.owner in (select username from users_list)
group by dba_segments.segment_name,dba_segments.owner;

exit
