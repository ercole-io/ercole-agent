-- Copyright (c) 2019 Sorint.lab S.p.A.

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

--column "TOTMB" format 9,999,990
--column "TBMB" format 9,999,990
--column "INDMB" format 9,999,990
--column "LOBMB" format 9,999,990
column "Nome_Acronimo" for a8
column "DB_Name" for a10

with 
lista_users as
(select username from dba_users where username not in ('SYS','AUDSYS','SYSTEM','SYSBACKUP','SYSDG','SYSKM','OUTLN','GSMADMIN_INTERNAL',
'GSMUSER','DIP','XS$NULL','ORACLE_OCM','DBSNMP','APPQOSSYS','ANONYMOUS','XDB',
'GSMCATUSER','WMSYS','OJVMSYS','CTXSYS','ORDDATA','ORDSYS','ORDPLUGINS',
'SI_INFORMTN_SCHEMA','MDSYS','OLAPSYS','MDDATA','SPATIAL_WFS_ADMIN_USR',
'SPATIAL_CSW_ADMIN_USR','LBACSYS','APEX_040200','APEX_PUBLIC_USER','FLOWS_FILES',
'DVSYS','DVF','SCOTT','EXFSYS','XS\$NULL','CMDB_USR','CPFI0_APPPRD','SYSMAN','SYSADMIN',
'MGMT_VIEW','DMSYS','WCCREPUSER','WCCUSER','WFADMIN','OWF_MGR','KWALKER','BLEWIS','CDOUGLAS','SPIERSON')
and username not like 'APEX%'
and username not like 'RMAN%'
and username not like 'U0%'
and username not like 'PATROL%'
and username not like 'OWB%'
and username not like '%WRPRD'
and username not like '%ROPRD'
and username not like '%GESTIONE'
and username not like 'FLOWS%'
and username not like 'SQLTX%'
and username not like 'DBA_%'
--and AUTHENTICATION_TYPE not like '%EXTERNAL%'
),
tbmb as
(select owner,round(sum(bytes/1024/1024)) as "A" from dba_segments where segment_type like 'TABLE%' group by owner),
indmb as
(select owner,round(sum(bytes/1024/1024)) as "B" from dba_segments where segment_type like 'INDEX%' group by owner),
lobmb as
(select owner,round(sum(bytes/1024/1024)) as "C" from dba_segments where segment_type like 'LOB%' group by owner)
select 
	   (select host_name from v$instance) as Hostname,
           (select value from v$parameter where name='db_name') as Nome_DB,
           (select db_unique_name from v$database) as DB_Unique_name,
	   u.username,
	   nvl(round(sum(s.bytes/1024/1024)),0) as "TOTMB",
	   nvl(t.A,0) as "TBMB",
	   nvl(i.B,0) as "INDMB",
	   nvl(l.C,0) as "LOBMB"
	   from 
	   lista_users u 
	   left join dba_segments s on  u.username=s.owner
	   left join tbmb t on  u.username=t.owner
	   left join indmb i on  u.username=i.owner
	   left join lobmb l on  u.username=l.owner
	   group by u.username,t.a,i.b,l.c order by 1;

exit
