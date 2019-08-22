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

col RECOMMENDATIONS for a300
col SEGMENT_OWNER for a60
col TABLESPACE_NAME like SEGMENT_OWNER
col RECL_GB for 9999999
col HOSTNAME for a255
col SEGMENT_NAME for a82
col SEGMENT_TYPE like SEGMENT_OWNER
col PARTITION_NAME like SEGMENT_OWNER

set autoprint on
set lines 900 pages 50000
SET FEEDBACK OFF
set colsep "|||"
VARIABLE x refcursor;
declare
  v_count NUMBER;
  v_count_res NUMBER;
  v_tot_db NUMBER;
  EDITION varchar2(80);
BEGIN

select (case when UPPER(banner) like '%EXTREME%' then 'EXTREME' when UPPER(banner) like '%ENTERPRISE%' then 'ENTERPRISE' else 'STANDARD' end) into EDITION from v$version where rownum=1;


IF ( EDITION ='STANDARD' ) THEN
   open :x for select
   (select HOST_NAME from v$instance) as HOSTNAME,
   (select name from v$database) as DB_NAME,
   (select '-' from dual) as SEGMENT_OWNER,
   (select '-' from dual) as SEGMENT_NAME,
   (select '-' from dual) as SEGMENT_TYPE,
   (select '-' from dual) as PARTITION_NAME,
   (select '-' from dual) as RECL_GB,
   (select '-' from dual) as RECOMMENDATIONS from dual where 1=0;

ELSE

select count(*) into v_count
from DBA_AUTOTASK_TASK
where CLIENT_NAME='auto space advisor'
and STATUS='ENABLED'
and trunc(LAST_GOOD_DATE) >= trunc(sysdate - 7);

IF(v_count > 0 ) THEN
	SELECT sum(bytes/1024/1024/1024) INTO v_tot_db FROM dba_segments;
	select count(*) into v_count_res from TABLE(dbms_space.asa_recommendations())
	where segment_owner not in ('ANONYMOUS','APEX_030200','APEX_040000','APEX_SSO','APPQOSSYS','CTXSYS','DBSNMP','DIP','EXFSYS','FLOWS_FILES','MDSYS','OLAPSYS','ORACLE_OCM','ORDDATA','ORDPLUGINS','ORDSYS','OUTLN','OWBSYS') 
  and segment_owner not in ('SI_INFORMTN_SCHEMA','SQLTXADMIN','SQLTXPLAIN','SYS','SYSMAN','SYSTEM','TRCANLZR','WMSYS','XDB','XS$NULL','PERFSTAT','STDBYPERF','MGDSYS','OJVMSYS')
  and round((100*(RECLAIMABLE_SPACE/1024/1024/1024))/v_tot_db,2) >= 0.1;
 
   IF(v_count_res > 0 ) THEN
     open :x for SELECT 
                 (select HOST_NAME from v$instance) as HOSTNAME,
                 (select name from v$database) as DB_NAME,
                 SEGMENT_OWNER,
                 SEGMENT_NAME,
                 SEGMENT_TYPE,
                 PARTITION_NAME,
                 decode(round(RECLAIMABLE_SPACE/1024/1024/1024,0),'0','<1',round(RECLAIMABLE_SPACE/1024/1024/1024,0)) as RECL_GB, 
                 RECOMMENDATIONS
                 FROM TABLE(dbms_space.asa_recommendations())
                 where segment_owner not in ('ANONYMOUS','APEX_030200','APEX_040000','APEX_SSO','APPQOSSYS','CTXSYS','DBSNMP','DIP','EXFSYS','FLOWS_FILES','MDSYS','OLAPSYS','ORACLE_OCM','ORDDATA','ORDPLUGINS','ORDSYS','OUTLN','OWBSYS') 
                 and segment_owner not in ('SI_INFORMTN_SCHEMA','SQLTXADMIN','SQLTXPLAIN','SYS','SYSMAN','SYSTEM','TRCANLZR','WMSYS','XDB','XS$NULL','PERFSTAT','STDBYPERF','MGDSYS','OJVMSYS')
and round((100*(RECLAIMABLE_SPACE/1024/1024/1024))/v_tot_db,2) >= 0.1
                 order by reclaimable_space desc;
        ELSE
		   open :x for select 
                      (select HOST_NAME from v$instance) as HOSTNAME, 
                      (select name from v$database) as DB_NAME,
                      (select '-' from dual) as SEGMENT_OWNER,
                      (select '-' from dual) as SEGMENT_NAME,
                      (select '-' from dual) as SEGMENT_TYPE,
                      (select '-' from dual) as PARTITION_NAME,
                      (select '-' from dual) as RECL_GB,
                      (select '-' from dual) as RECOMMENDATIONS from dual where 1=0;
	END IF;

  ELSE
   open :x for select 
   (select HOST_NAME from v$instance) as HOSTNAME, 
   (select name from v$database) as DB_NAME,
   (select '-' from dual) as SEGMENT_OWNER,
   (select '-' from dual) as SEGMENT_NAME,
   (select '-' from dual) as SEGMENT_TYPE,
   (select '-' from dual) as PARTITION_NAME,
   (select '-' from dual) as RECL_GB,
   (select '-' from dual) as RECOMMENDATIONS from dual where 1=0;

END IF;		
END IF;		
END;
/
exit
