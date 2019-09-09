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

set lines 8000 pages 0 feedback off verify off
set colsep "|||"
set serverout on
set autoprint on
col HOST_NAME for a30
col TASK_NAME for a15 
col "BENEFIT_%" for 999.99
col FINDING for a600 
col RECOMMANDETION for a30
col ACTION like FINDING

VARIABLE x refcursor;
define num_days = 7;
declare
  v_dbid NUMBER;
  v_inst_num NUMBER;
  v_begin_snap NUMBER;
  v_end_snap NUMBER;
  v_count_usage NUMBER;
  v_count_benefit NUMBER;
  v_tname VARCHAR2(30);
  EDITION varchar2(80); 
BEGIN

select (case when UPPER(banner) like '%EXTREME%' then 'EXTREME' when UPPER(banner) like '%ENTERPRISE%' then 'ENTERPRISE' else 'STANDARD' end) into EDITION from v$version where rownum=1;

IF ( EDITION = 'STANDARD' ) THEN
    open :x for select
                      (select HOST_NAME from v$instance) as HOST_NAME,
                      (select instance_name from v$instance) as INSTANCE_NAME,
                      (select 'N/A' from dual) as FINDING,
                      (select 'N/A' from dual) as RECOMMANDETION,
                      (select 'N/A' from dual) as ACTION,
                      (select 'N/A' from dual) as "BENEFIT_%" from dual where 1=0;
ELSE
with strtime as
(
  select  max(startup_time) as data, 
          max(count(*)) 
  from dba_hist_snapshot 
  where instance_number=(select instance_number from v$instance)
  group by startup_time
)
SELECT NVL(MIN(s.snap_id),0), 
       NVL(MAX(s.snap_id),0) INTO v_begin_snap,
       v_end_snap
FROM dba_hist_snapshot s, 
     strtime st
where s.startup_time = st.data 
and trunc(s.END_INTERVAL_TIME) >= trunc(sysdate - &&num_days) ;

  IF (v_begin_snap=0) OR (v_end_snap=0) THEN
    open :x for select 
                      (select HOST_NAME from v$instance) as HOST_NAME, 
                      (select instance_name from v$instance) as INSTANCE_NAME,
                      (select 'N/A' from dual) as FINDING,
                      (select 'N/A' from dual) as RECOMMANDETION,
                      (select 'N/A' from dual) as ACTION,
                      (select 'N/A' from dual) as "BENEFIT_%" from dual where 1=0; 
   ELSE
     SELECT dbid INTO v_dbid FROM v$database;
     SELECT instance_number INTO v_inst_num FROM v$instance;
     SELECT 'ERCOLE_TASK' into v_tname from dual; 
     /* Check if ADDM is enable */
     select max(DETECTED_USAGES) INTO v_count_usage from dba_feature_usage_statistics where name ='AWR Report' and dbid=v_dbid;
  
     IF (v_count_usage = 0) THEN
       open :x for select 
                        (select HOST_NAME from v$instance) as HOST_NAME, 
                        (select instance_name from v$instance) as INSTANCE_NAME,
                        (select 'N/A' from dual) as FINDING,
                        (select 'N/A' from dual) as RECOMMANDETION,
                        (select 'N/A' from dual) as ACTION,
                        (select 'N/A' from dual) as "BENEFIT_%" from dual where 1=0; 
      ELSE
        /* Create ADDM TASK for instance */
        dbms_addm.analyze_inst(v_tname,v_begin_snap,v_end_snap,v_inst_num,v_dbid);
        with 
          dbtime as
          (
            /* Extract DB_TIME for Task*/
            select f.TASK_NAME,sum(unique(i.DATABASE_TIME)) as DB_TIME
            from dba_addm_instances i, DBA_ADDM_FINDINGS f
            where i.TASK_ID=f.TASK_ID
            group by f.TASK_NAME
          )
          SELECT count(*) INTO v_count_benefit from DBA_ADDM_FINDINGS f, DBA_ADVISOR_RECOMMENDATIONS r, DBA_ADVISOR_ACTIONS a, dbtime d
          where f.TASK_ID=r.TASK_ID
          and a.TASK_ID=r.TASK_ID
          and f.TASK_NAME=v_tname
          and d.TASK_NAME=f.TASK_NAME
          and r.REC_ID=a.REC_ID
          and f.FINDING_ID=r.FINDING_ID
          and ROUND((100*r.BENEFIT)/d.DB_TIME,2) >=20;
  
        IF (v_count_benefit = 0) THEN
          open :x for select 
                        (select HOST_NAME from v$instance) as HOST_NAME, 
                        (select instance_name from v$instance) as INSTANCE_NAME,
                        (select 'N/A' from dual) as FINDING,
                        (select 'N/A' from dual) as RECOMMANDETION,
                        (select 'N/A' from dual) as ACTION,
                        (select 'N/A' from dual) as "BENEFIT_%" from dual where 1=0; 
         ELSE
           /*Extract REPORT info */
           open :x for with 
             dbtime as
             (
               /* Extract DB_TIME for Task*/
               select f.TASK_NAME,sum(unique(i.DATABASE_TIME)) as DB_TIME
               from dba_addm_instances i, DBA_ADDM_FINDINGS f
               where i.TASK_ID=f.TASK_ID
               group by f.TASK_NAME
             ),
             server as
             (
               select host_name,instance_name from v$instance
             )
             /* Extract Finding, Recommendation, Action and Benefit %*/
             select s.host_name,s.instance_name, f.MESSAGE as FINDING, r.TYPE||CHR(10)||r.BENEFIT_TYPE as RECOMMANDETION, a.MESSAGE as ACTION, ROUND((100*r.BENEFIT)/d.DB_TIME,2) as "BENEFIT_%"
             from DBA_ADDM_FINDINGS f, DBA_ADVISOR_RECOMMENDATIONS r, DBA_ADVISOR_ACTIONS a, dbtime d, server s
             where f.TASK_ID=r.TASK_ID
             and a.TASK_ID=r.TASK_ID
             and f.TASK_NAME=v_tname
             and d.TASK_NAME=f.TASK_NAME
             and r.REC_ID=a.REC_ID
             and f.FINDING_ID=r.FINDING_ID
             and ROUND((100*r.BENEFIT)/d.DB_TIME,2) >=20
             order by 6 desc;
        END IF;      
     END IF;     
  END IF;
END IF;
END;
/

/* Delete ADDM TASK for instance */
exec dbms_addm.delete('ERCOLE_TASK');

exit
