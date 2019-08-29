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

set lines 16000 pages 0 feedback off verify off
set colsep "|||"
col owner for a30
col Nome_Acronimo for a8
col segment_name for a60

VARIABLE dbid NUMBER;
VARIABLE inst_num NUMBER;
VARIABLE bid NUMBER;
VARIABLE eid NUMBER;
VARIABLE elapsed varchar2(100);
VARIABLE dbtime  varchar2(100);
VARIABLE count_usage NUMBER;
VARIABLE result varchar2(100);
VARIABLE esclusion varchar2(100);
BEGIN
with strtime as
(select max(startup_time) as data, max(count(*)) from dba_hist_snapshot where instance_number=(select instance_number from v$instance) group by startup_time)
SELECT MIN (s.snap_id), MAX (s.snap_id) INTO :bid,:eid
FROM dba_hist_snapshot s, strtime st
where s.startup_time=st.data;
SELECT dbid INTO :dbid FROM v$database;
SELECT instance_number INTO :inst_num FROM v$instance;
select max(DETECTED_USAGES) INTO :count_usage from dba_feature_usage_statistics where name ='AWR Report' and dbid=:dbid;
-- Esclusione DATABASE 
select name INTO :esclusion from v$database;
IF(:esclusion ='CRPL0ONP' or :esclusion ='GOLD0ONP') THEN
	:count_usage := 0;
END IF;

IF (:count_usage > 0) THEN
   	SELECT replace(replace(replace(output,'Elapsed:',''),chr(32), ''),'(mins)','')  into :elapsed FROM TABLE (DBMS_WORKLOAD_REPOSITORY.awr_report_text (:dbid, :inst_num, :bid, :eid, 0)) where rownum <2 and output like '%Elapsed: %';
	SELECT replace(replace(replace(output,'DB Time:',''),chr(32), ''),'(mins)','') into :dbtime FROM TABLE (DBMS_WORKLOAD_REPOSITORY.awr_report_text (:dbid, :inst_num, :bid, :eid, 0)) where rownum <2 and output like '%DB Time: %';
select to_char(round(((to_number(:dbtime,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''')/to_number(:elapsed,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,'''))),4),'99990') into :result 
from dual;

IF (:result = 0) THEN
select '1' into :result from dual;
END IF;

ELSE
   	 SELECT 'N/A' into :elapsed from dual;
   	 SELECT 'N/A' into :dbtime from dual;
   	 SELECT 'N/A' into :result from dual;
END IF;
END;
/


select 
(select value from v$parameter where name='db_name') as Nome_DB,
(select db_unique_name from v$database) as DB_Unique_name,
(select instance_number from v$instance) as Instance_number,
(select status from v$instance) as DB_Status,
((SELECT version FROM V$INSTANCE)||(select (case when UPPER(banner) like '%EXTREME%' then ' Extreme Edition' when UPPER(banner) like '%ENTERPRISE%' then ' Enterprise Edition' else ' Standard Edition' end) from v$version where rownum=1)) as Versione,
(SELECT platform_name  FROM V$database) as platform,
(SELECT log_mode  FROM V$database) as archive,
(select value from nls_database_parameters where parameter='NLS_CHARACTERSET') as Charset,
(select value from nls_database_parameters where parameter='NLS_NCHAR_CHARACTERSET') as NCharset,
(select value from v$parameter where name='db_block_size') as Blocksize,
(select value from v$parameter where name='cpu_count') as Cpu_count,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS='',.'''),',') from v$parameter where name='sga_target')  as Sga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS='',.'''),',') from v$parameter where name='pga_aggregate_target') as Pga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS='',.'''),',') from v$parameter where name='memory_target') as Pga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS='',.'''),',') from v$parameter where name='sga_max_size') as sga_max_size,
(select round(sum(bytes/1024/1024/1024)) from dba_segments) as Alloc,
((select round(sum(bytes/1024/1024/1024)) from dba_data_files)+(select round(sum(bytes/1024/1024/1024)) from dba_temp_files)+(select round(sum(bytes/1024/1024/1024)) from v$log)),
((select round(sum(decode(autoextensible,'NO',bytes/1024/1024/1024,'YES',maxbytes/1024/1024/1024))) from dba_data_files)+(select round(sum(bytes/1024/1024/1024)) from dba_temp_files)+(select round(sum(bytes/1024/1024/1024)) from v$log)),
(select 
	case when (select :elapsed from dual) != 'N/A' 
	then 
		to_number(:elapsed,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''') 
	else 0 end as "elapsed" from dual),
(select case when (select :dbtime from dual) != 'N/A' then 
to_number(:dbtime,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''') else 0 end as "dbtime" from dual),
(select :result from dual),
(select 
	case when (select count(*) from dba_data_files where file_name like '+%') > 0 
		then 'Y' 
		else 'N' end as "ASM" from dual ),
	case when ( select count(*) from V_$DATAGUARD_CONFIG) > 1 
		then 'Y' 
		else 'N' end  as "Dataguard" 
from dual;

exit
