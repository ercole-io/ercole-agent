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

set lines 32767 pages 0 feedback off verify off
set colsep "|||"
col owner for a30
col Nome_Acronimo for a8
col segment_name for a60

select 
(select value from v$parameter where name='db_name') as Nome_DB,
(select db_unique_name from v$database) as DB_Unique_name,
(select instance_number from v$instance) as Instance_number,
(select status from v$instance) as DB_Status,
((SELECT version FROM V$INSTANCE)||(select (case when UPPER(banner) like '%EXTREME%' then ' Extreme Edition' when UPPER(banner) like '%ENTERPRISE%' then ' Enterprise Edition' else ' Standard Edition' end) from v$version where rownum=1)) as Versione,
(SELECT platform_name  FROM V$database) as platform,
(SELECT log_mode  FROM V$database) as archive,
(select value from v$nls_parameters where parameter='NLS_CHARACTERSET') as Charset,
(select value from v$nls_parameters where parameter='NLS_NCHAR_CHARACTERSET') as NCharset,
(select value from v$parameter where name='db_block_size') as Blocksize,
(select value from v$parameter where name='cpu_count') as Cpu_count,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',') from v$parameter where name='sga_target')  as Sga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',') from v$parameter where name='pga_aggregate_target') as Pga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',') from v$parameter where name='memory_target') as Pga_Target,
(select rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',') from v$parameter where name='sga_max_size') as sga_max_size,
--(select round(sum(bytes/1024/1024/1024)) from dba_segments),
'0',
--((select round(sum(bytes/1024/1024/1024)) from v$datafile)+(select round(sum(bytes/1024/1024/1024)) from v$tempfile)+(select round(sum(bytes/1024/1024/1024)) from v$log)),
'0',
--((select round(sum(decode(autoextensible,'NO',bytes/1024/1024/1024,'YES',maxbytes/1024/1024/1024))) from v$datafile)+(select round(sum(bytes/1024/1024/1024)) from v$tempfile)+(select round(sum(bytes/1024/1024/1024)) from v$log)),
'0',
--(SELECT replace(replace(output,'Elapsed:',''),chr(32), '') FROM TABLE (DBMS_WORKLOAD_REPOSITORY.awr_report_text (:dbid, :inst_num, :bid, :eid, 0)) where rownum <20 and output like '%Elapsed: %'),
--(SELECT replace(replace(output,'DB Time:',''),chr(32), '') FROM TABLE (DBMS_WORKLOAD_REPOSITORY.awr_report_text (:dbid, :inst_num, :bid, :eid, 0)) where rownum <20 and output like '%DB Time: %'),
--:elapsed,:dbtime,(select :result from dual),
'0','0','0','0',
(select case when (select count(*) from v$datafile where name like '+%') > 0 then 'Y' else 'N' end as "ASM" from dual ),
case when ( select count(*) from V$DATAGUARD_CONFIG) > 1 then 'Y' else 'N' end  as "Dataguard"
from dual;
exit
