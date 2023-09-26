-- Copyright (c) 2023 Sorint.lab S.p.A.

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

set feedback off heading off pages 0 serverout on verify off lines 1234 timing off

column 1 new_value 1 noprint
select '' "1" from dual where rownum = 0;
--If the script has been called without argument (should be in case PDB name) it has to be executed on the CDB/TRADITIONAL DB
define param = &1 '-1'

begin
        if '&param' <> '-1' then
                execute immediate 'alter session set container=&param';
        end if;
end;
/

--If version >= 11.1 COMPRESS_FOR column exists in tables dba_tables and dba_tab_partitions, it doesn't exists in previous versions
--it is simulated a where to return 0 rows
define COMPRESS_FOR = 'COMPRESS_FOR'
col :COMPRESS_FOR_ new_val COMPRESS_FOR noprint
variable COMPRESS_FOR_ varchar2(30) 

DECLARE 
DB_VERSION number;
begin
	DB_VERSION := dbms_db_version.version + (dbms_db_version.release / 10);
	if DB_VERSION >= 11.1 then
		:COMPRESS_FOR_ := 'COMPRESS_FOR';
	else
		:COMPRESS_FOR_ := '''BASIC''';
	end if;	
end;
/

select :COMPRESS_FOR_ from dual;

set colsep "|||"

select 'PLSQL LINES', nvl(count(line),0) line_total from dba_source where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
union all
select 'PARTITIONED TABLES', count(*) from dba_tables 
where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
and PARTITIONED='YES'
union all
select 'PARTITIONED INDEXES', count(*) from dba_indexes 
where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
and PARTITIONED='YES'
union all
select TEMP.description,nvl(dba_objects.num_objs, 0) from  
(select object_type,count(*) num_objs from dba_objects
where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
and object_type in ('SEQUENCE','PROCEDURE','FUNCTION','TRIGGER')
group by object_type) dba_objects
RIGHT JOIN 
(select * from 
(select 'SEQUENCE' object_type,'SEQUENCES' description from dual
union all 
select 'PROCEDURE' object_type,'PROCEDURES' description from dual
union all
select 'FUNCTION' object_type,'FUNCTIONS' description from dual
union all
select 'TRIGGER' object_type,'TRIGGERS' description from dual)) TEMP
on TEMP.object_type=dba_objects.object_type
union all
select 'HCC COMPRESSED TABLES',sum(num_objs) from
(select count(*) num_objs from dba_tables
where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
and compression='ENABLED'
and &&COMPRESS_FOR not in ('BASIC','ADVANCED')
union all  
select count(distinct table_owner||'.'||table_name) num_objs from dba_tab_partitions
where table_owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and table_owner not like 'APEX%'
and compression='ENABLED'
and &&COMPRESS_FOR not in ('BASIC','ADVANCED'))
union all
select 'MVIEWS REWRITE ENABLED',count(*) from dba_mviews
where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
and rewrite_enabled='Y'
union all
select 'VDP POLICIES',count(*) num_objs from dba_policies
where object_owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and object_owner not like 'APEX%'
and enable='YES';

select owner,type,nvl(count(line),0) line_total from dba_source where owner not in ('RDSADMIN','SYS','SYSTEM','XDB','X$NULL','SYSMAN','SYSBACKUP','SPATIAL%','SI_INFORMTN_SCHEMA','ORDSYS','ORDDATA','OUTLN','OJVMSYS','OWB%','MDSYS','ORACLE_OCM','OLAPSYS','APEX_040200','APEX_030200','AUDSYS','APPQOSSYS','CTXSYS','DBSNMP','DVSYS','FLOWS_FILES','GSMADMIN_INTERNAL','LBACSYS','ORDPLUGINS','OWBSYS','OWBSYS_AUDIT','WMSYS','REMOTE_SCHEDULER_AGENT','DVF','DBSFWUSER','SYSDG','SYSKM','XS$NULL','GSMUSER','DIP','ANONYMOUS','GSMCATUSER','MDDATA','APEX_PUBLIC_USER','SYSRAC','SYS$UMF','GGSYS','SPATIAL_CSW_ADMIN_USR') 
and owner not like 'APEX%'
group by owner,type 
order by owner,type;

exit