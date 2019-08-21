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

set feedback off

VARIABLE RETENTION varchar2(200)
VARIABLE BCK number;

begin
select count(*) into :BCK from  v$rman_configuration where name ='RETENTION POLICY';

IF (:BCK > 0 ) THEN
  SELECT 
    CASE 
      WHEN value LIKE 'TO RECOVERY WINDOW%'  THEN   regexp_replace(value, '[^0-9]', '')|| ' DAYS'
      WHEN value LIKE 'TO REDUNDANCY%'  THEN   regexp_replace(value, '[^0-9]', '')|| ' NUMBERS'
    end into :RETENTION
  from  v$rman_configuration where name ='RETENTION POLICY';
ELSE
  select '1 NUMBERS' into :RETENTION from dual;
END IF;  
end;
/


alter session set nls_date_format='dd-mm-yyyy hh24:mi';

set lines 200 pages 0
set colsep "|||"
col WEEK_DAYS for a70
col RETENTION for a20

with backup as
(
  select 
    j.start_time start_time,
    j.end_time end_time,
    decode(j.input_type,'DB INCR',decode(i0,0,'Incr Lvl 1','Incr Lvl 0'),initcap(j.input_type)) backup_type,
    round(sum (j.output_bytes/1024/1024)) AVG_BCK_SIZE_MB
  from V$RMAN_BACKUP_JOB_DETAILS j
  left outer join 
      (select d.session_recid,
            d.session_stamp,
            sum(case when d.backup_type = 'D' then d.pieces else 0 end) D,
            sum(case when (d.backup_type||d.incremental_level = 'I0' or d.incremental_level = 0) then d.pieces else 0 end) I0,
            sum(case when (d.backup_type||d.incremental_level = 'I1' or d.incremental_level = 1) then d.pieces else 0 end) I1,
            sum(case when d.backup_type = 'L' then d.pieces else 0 end) L
      from
      V$BACKUP_SET_DETAILS d join V$BACKUP_SET s on s.set_stamp = d.set_stamp and s.set_count = d.set_count
      where s.input_file_scan_only = 'NO'
      group by d.session_recid, d.session_stamp
      ) x on x.session_recid = j.session_recid and x.session_stamp = j.session_stamp
  left outer join 
      (select o.session_recid, o.session_stamp, min(inst_id) inst_id
       from GV$RMAN_OUTPUT o
       group by o.session_recid, o.session_stamp
      ) ro on ro.session_recid = j.session_recid and ro.session_stamp = j.session_stamp
  --where j.start_time > trunc(next_day(sysdate-30,'SATURDAY'))
  where j.start_time > trunc(sysdate-30)
  and j.status not in ('FAILED')
  group by start_time,end_time, decode(j.input_type,'DB INCR',decode(i0,0,'Incr Lvl 1','Incr Lvl 0'),initcap(j.input_type))
  order by j.start_time 
),
output as
(
select backup_type,
       /* Round to nearest quarter */
       to_char(TRUNC(start_time)+(ROUND((start_time - TRUNC (start_time))* 96)/ 96),'HH24:MI') as HOUR,
       count(*) as N_TIMES,
       to_char(start_time, 'FmDay','nls_date_language=English') as WEEK_DAY,
       round(avg(AVG_BCK_SIZE_MB)/1024) as AVG_BCK_SIZE_GB
from backup  
where backup_type in ('Db Full','Incr Lvl 0','Incr Lvl 1','Archivelog','Recvr Area','Datafile Full','Datafile Incr') 
group by to_char(TRUNC(start_time)+(ROUND((start_time - TRUNC (start_time))* 96)/ 96),'HH24:MI'), backup_type,to_char(start_time, 'FmDay','nls_date_language=English')
),
appo_multi_bck as
(
select distinct BACKUP_TYPE,HOUR
from output
where N_TIMES > 1
order by 1
),
one_time_bck as
(
select o.backup_type,
       o.hour
from output o 
left join appo_multi_bck a on o.backup_type = a.backup_type and o.hour = a.hour
where o.backup_type not in (select a.BACKUP_TYPE from appo_multi_bck a)
),
hour as
(
select backup_type,hour
from appo_multi_bck
union all
select backup_type,hour
from one_time_bck 
)
select o.backup_type, 
       o.HOUR,
       rtrim(xmlagg(xmlelement(e, o.WEEK_DAY, ',')).extract('//text()').getclobval(), ', ') as WEEK_DAYS,
       round(avg(o.AVG_BCK_SIZE_GB)) as AVG_BCK_SIZE_GB,
       :RETENTION as RETENTION
from output o 
right outer join hour h on o.backup_type = h.backup_type and o.hour = h.hour
group by o.backup_type,o.HOUR
order by 1,2;
exit

/*
Esempio output:

BACKUP_TYPE   HOUR  WEEK_DAYS                                                              AVG_BCK_SIZE_GB RETENTION
------------- ----- ---------------------------------------------------------------------- --------------- --------------------
Archivelog    00:00 Tuesday,Saturday,Monday,Friday,Sunday,Thursday,Wednesday                             8 150 DAYS
Archivelog    06:00 Monday,Saturday,Friday,Thursday,Sunday,Tuesday,Wednesday                             6 150 DAYS
Archivelog    12:00 Tuesday,Sunday,Thursday,Wednesday,Friday,Saturday,Monday                             3 150 DAYS
Archivelog    18:00 Monday,Thursday,Tuesday,Wednesday,Sunday,Saturday,Friday                             4 150 DAYS
Incr Lvl 0    01:00 Sunday                                                                             627 150 DAYS
Incr Lvl 1    01:00 Friday,Monday,Thursday,Saturday,Tuesday,Wednesday                                   38 150 DAYS
*/