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

VARIABLE dbid NUMBER;
VARIABLE inst_num NUMBER;
VARIABLE bid NUMBER;
VARIABLE eid NUMBER;
VARIABLE elapsed varchar2(100);
VARIABLE dbtime  varchar2(100);
VARIABLE cputime  varchar2(100);
VARIABLE count_usage NUMBER;
VARIABLE result varchar2(100);
VARIABLE esclusion varchar2(100);

BEGIN
WITH strtime AS
  (SELECT max(startup_time) AS DATA,
          max(count(*))
   FROM dba_hist_snapshot
   WHERE BEGIN_INTERVAL_TIME > trunc(sysdate-nvl('&&1',30))
     AND instance_number=
       (SELECT instance_number
        FROM v$instance)
   GROUP BY startup_time)
SELECT MIN (s.snap_id), MAX (s.snap_id) INTO :bid,
                                             :eid
FROM dba_hist_snapshot s,
     strtime st
WHERE s.startup_time=st.DATA
  AND s.BEGIN_INTERVAL_TIME > trunc(sysdate-nvl('&&1',30));


SELECT dbid INTO :dbid FROM v$database;
SELECT instance_number INTO :inst_num FROM v$instance;


SELECT max(DETECTED_USAGES) INTO :count_usage
FROM dba_feature_usage_statistics
WHERE name ='AWR Report'
  AND dbid=:dbid;

-- Esclusione DATABASE 
select name INTO :esclusion from v$database;
IF(:esclusion ='DB01' or :esclusion ='DB02') THEN
	:count_usage := 0;
END IF;

IF (:count_usage > 0) THEN
   	WITH awrr AS
	  (SELECT *
	   FROM TABLE (DBMS_WORKLOAD_REPOSITORY.awr_report_text (:dbid, :inst_num, :bid, :eid, 0))
	   WHERE rownum <100)
	SELECT
	  (SELECT replace(replace(replace(OUTPUT,'Elapsed:',''),chr(32), ''),'(mins)','')
	   FROM awrr
	   WHERE rownum <2
	     AND OUTPUT LIKE '%Elapsed: %') AS a,
	  (SELECT replace(replace(replace(OUTPUT,'DB Time:',''),chr(32), ''),'(mins)','')
	   FROM awrr
	   WHERE rownum <2
	     AND OUTPUT LIKE '%DB Time: %') AS b,
	  (SELECT REGEXP_SUBSTR(replace(replace(OUTPUT,'DB CPU(s):',''),chr(32), '|'),'[^|]+',1,1)
	   FROM awrr
	   WHERE rownum <2
	     AND OUTPUT LIKE '%DB CPU(s): %') AS c INTO :elapsed,
	                                                :dbtime,
	                                                :cputime
	FROM awrr
	WHERE rownum <2;

	select to_char(round(((to_number(:dbtime,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''')/to_number(:elapsed,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,'''))),4),'99990') into :result 
	from dual;

IF (:result = 0) THEN
select '1' into :result from dual;
END IF;
ELSE
   	 SELECT 'N/A' into :elapsed from dual;
   	 SELECT 'N/A' into :dbtime from dual;
   	 SELECT 'N/A' into :result from dual;
   	 SELECT 'N/A' into :cputime from dual;
END IF;
END;
/

SELECT
  (SELECT value
   FROM v$parameter
   WHERE name='db_name') AS Nome_DB,
  (SELECT db_unique_name
   FROM v$database) AS DB_Unique_name,
  (SELECT instance_number
   FROM v$instance) AS Instance_number,
  (SELECT status
   FROM v$instance) AS DB_Status, (
                                     (SELECT VERSION
                                      FROM V$INSTANCE)||
                                     (SELECT (CASE WHEN UPPER(banner) LIKE '%EXTREME%' THEN ' Extreme Edition' WHEN UPPER(banner) LIKE '%ENTERPRISE%' THEN ' Enterprise Edition' ELSE ' Standard Edition' END)
                                      FROM v$version
                                      WHERE rownum=1)) AS Versione,
  (SELECT platform_name
   FROM V$database) AS platform,
  (SELECT log_mode
   FROM V$database) AS archive,
  (SELECT value
   FROM nls_database_parameters
   WHERE PARAMETER='NLS_CHARACTERSET') AS Charset,
  (SELECT value
   FROM nls_database_parameters
   WHERE PARAMETER='NLS_NCHAR_CHARACTERSET') AS NCharset,
  (SELECT value
   FROM v$parameter
   WHERE name='db_block_size') AS Blocksize,
  (SELECT value
   FROM v$parameter
   WHERE name='cpu_count') AS Cpu_count,
  (SELECT rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',')
   FROM v$parameter
   WHERE name='sga_target') AS Sga_Target,
  (SELECT rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',')
   FROM v$parameter
   WHERE name='pga_aggregate_target') AS Pga_Target,
  (SELECT rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',')
   FROM v$parameter
   WHERE name='memory_target') AS Pga_Target,
  (SELECT rtrim(to_char(value/1024/1024/1024, 'FM9G999G999D999', 'NLS_NUMERIC_CHARACTERS=''.,'''),',')
   FROM v$parameter
   WHERE name='sga_max_size') AS sga_max_size,
  (SELECT round(sum(bytes/1024/1024/1024))
   FROM dba_segments) AS Alloc, (
                                   (SELECT round(sum(bytes/1024/1024/1024))
                                    FROM dba_data_files)+
                                   (SELECT round(sum(bytes/1024/1024/1024))
                                    FROM dba_temp_files)+
                                   (SELECT round(sum(bytes/1024/1024/1024))
                                    FROM v$log)), (
                                                     (SELECT round(sum(decode(autoextensible,'NO',bytes/1024/1024/1024,'YES',maxbytes/1024/1024/1024)))
                                                      FROM dba_data_files)+
                                                     (SELECT round(sum(bytes/1024/1024/1024))
                                                      FROM dba_temp_files)+
                                                     (SELECT round(sum(bytes/1024/1024/1024))
                                                      FROM v$log)),
  (SELECT CASE
              WHEN
                     (SELECT :elapsed
                      FROM dual) != 'N/A' THEN to_number(:elapsed,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''')
              ELSE 0
          END AS "elapsed"
   FROM dual),
  (SELECT CASE
              WHEN
                     (SELECT :dbtime
                      FROM dual) != 'N/A' THEN to_number(:dbtime,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''')
              ELSE 0
          END AS "dbtime"
   FROM dual),
  (SELECT CASE
              WHEN
                     (SELECT :cputime
                      FROM dual) != 'N/A' THEN to_number(:cputime,'999999.99',' NLS_NUMERIC_CHARACTERS = ''.,''')
              ELSE 0
          END AS "cputime"
   FROM dual),
  (SELECT :result
   FROM dual),
  (SELECT CASE
              WHEN
                     (SELECT count(*)
                      FROM dba_data_files
                      WHERE file_name LIKE '+%') > 0 THEN 'Y'
              ELSE 'N'
          END AS "ASM"
   FROM dual), CASE
                   WHEN
                          (SELECT count(*)
                           FROM V$DATAGUARD_CONFIG) > 1 THEN 'Y'
                   ELSE 'N'
               END AS "Dataguard"
FROM dual;

exit
