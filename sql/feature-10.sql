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

set feedback off pages 0 serverout on verify off lines 1234
VARIABLE TUNING varchar2(100);
VARIABLE DIAGNOSTICS varchar2(100);
VARIABLE LABELSECURITY varchar2(100);
VARIABLE COMPRESSION varchar2(100);
VARIABLE ANALYTICS varchar2(100);
VARIABLE TESTING varchar2(100);
VARIABLE OLAP varchar2(100);
VARIABLE VAULT varchar2(100);
VARIABLE PARTITIONING varchar2(100);
VARIABLE RAC varchar2(100);
VARIABLE SPATIAL varchar2(100);
VARIABLE GATEWAY varchar2(100);
--
VARIABLE WEBLOGICSERVER  varchar2(100);
VARIABLE SECURE varchar2(100);
VARIABLE ONE varchar2(100);
VARIABLE PATCH1 varchar2(100);
VARIABLE PATCH2 varchar2(100);
VARIABLE PILLAR varchar2(100);
VARIABLE MULTITENANT varchar2(100);
VARIABLE HW varchar2(100);
VARIABLE GOLDEN varchar2(100);
VARIABLE EXADATA varchar2(100);
VARIABLE MASKING varchar2(100);
VARIABLE MEMORY varchar2(100);
VARIABLE CONFIGURATION varchar2(100);
VARIABLE SECURITY varchar2(100);
VARIABLE MANAGEMENT varchar2(100);
VARIABLE GUARD varchar2(100);
VARIABLE RAC2 varchar2(100);
begin
select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Automatic Maintenance - SQL Tuning Advisor','Automatic SQL Tuning Advisor','Real-Time SQL Monitoring','Real-Time SQL Monitoring','SQL Access Advisor','SQL Monitoring and Tuning pages','SQL Profile','SQL Tuning Advisor',
		'SQL Tuning Set (user)','Tuning Pack') 
		and detected_usages > 0 and dbid=(select dbid from v$database) )) > 0 
	then 'Y' else '' end into :TUNING from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('AWR Repository','Automatic Database Diagnostic Monitor','Automatic Workload Repository','Diagnostic Pack') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :DIAGNOSTICS from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Label Security') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :LABELSECURITY from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Data Guard') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :COMPRESSION from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Data Mining') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :ANALYTICS from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Database Replay: Workload Capture') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :TESTING from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('OLAP - Analytic Workspaces','OLAP - Cubes') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :OLAP from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Oracle Database Vault') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :VAULT from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Partitioning (user)') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :PARTITIONING  from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Real Application Clusters (RAC)') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :RAC from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Transparent Gateway') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :GATEWAY from dual;

select case when 
	(select count(*) 
		from dba_feature_usage_statistics 
		where (name in('Spatial') 
		and detected_usages > 0 and dbid=(select dbid from v$database))) > 0 
	then 'Y' else '' end into :SPATIAL from dual;
end;
/

select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'WebLogic Server Management Pack Enterprise Edition:'||:WEBLOGICSERVER
else 'WebLogic Server Management Pack Enterprise Edition:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Tuning Pack:'||:TUNING
else 'Tuning Pack:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Spatial and Graph:'||:SPATIAL
else 'Spatial and Graph:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Secure Backup:'||:SECURE
else 'Secure Backup:' end from dual;
select 'Real Application Testing:'||:TESTING from dual;
select 'Real Application Clusters One Node:'||:ONE from dual;
select 'Real Application Clusters:'||:RAC from dual;
select 'RAC or RAC One Node:'||:RAC2 from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Provisioning and Patch Automation Pack for Database:'||:PATCH1
else 'Provisioning and Patch Automation Pack for Database:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Provisioning and Patch Automation Pack:'||:PATCH2
else 'Provisioning and Patch Automation Pack:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Pillar Storage:'||:PILLAR
else 'Pillar Storage:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Partitioning:'||:PARTITIONING
else 'Partitioning:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'OLAP:'||:OLAP
else 'OLAP:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Multitenant:'||:MULTITENANT
else 'Multitenant:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Label Security:'||:LABELSECURITY
else 'Label Security:' end from dual;
select
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'HW:'||:HW
else 'HW:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'GoldenGate:'||:GOLDEN
else 'GoldenGate:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Exadata:'||:EXADATA
else 'Exadata:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Diagnostics Pack:'||:DIAGNOSTICS
else 'Diagnostics Pack:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Database Vault:'||:VAULT
else 'Database Vault:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Database In-Memory:'||:MEMORY
else 'Database In-Memory:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Database Gateway:'||:GATEWAY
else 'Database Gateway:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Data Masking Pack:'||:MASKING
else 'Data Masking Pack:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Configuration Management Pack for Oracle Database:'||:CONFIGURATION
else 'Configuration Management Pack for Oracle Database:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Change Management Pack:'||:MANAGEMENT
else 'Change Management Pack:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Advanced Security:'||:SECURITY
else 'Advanced Security:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Advanced Compression:'||:COMPRESSION
else 'Advanced Compression:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Advanced Analytics:'||:ANALYTICS
else 'Advanced Analytics:' end from dual;
select 
case when (select UPPER(banner) from v$version where rownum=1) like '%EXTREME%' or (select UPPER(banner) from v$version where rownum=1) like '%ENTERPRISE%' then
'Active Data Guard:'||:GUARD
else 'Active Data Guard:' end from dual;
