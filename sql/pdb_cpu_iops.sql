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

set lines 8000 pages 0

define SNAPDAYS = 30;
define MONTH = 30;
define WEEKS_NUMBER = 4;
define WEEK = 7;
define MONTH_SOFT_LIMIT = 15;
define WEEK_SOFT_LIMIT = 4;

column bid new_value V_BID noprint
column eid new_value V_EID noprint
column dbid new_value V_DBID noprint
column con_id new_value V_CONID noprint
column instance_number new_value V_INSTANCE_NUMBER noprint
column snapdaysretrieved new_value V_SNAP_DAYS_RETRIEVED noprint
column diagnostic_pack_utilizations new_value V_diagnostic_pack_utilizations noprint

--Fetch the database dbid
SELECT dbid FROM v$database;

select nvl(sum(detected_usages),0) diagnostic_pack_utilizations
from DBA_FEATURE_USAGE_STATISTICS
where name in ('ADDM','AWR Baseline','AWR Baseline Template','AWR Report','Automatic Workload Repository','Baseline Adaptive Thresholds',
'Baseline Static Computations','Diagnostic Pack','EM Performance Page') and dbid=&V_DBID;

alter session set container=&1;

set serveroutput on

--Fetch the database dbid
SELECT instance_number FROM v$instance;

--Consering that DBA_HIST_CON_SYSMETRIC_SUMM and DBA_HIST_RSRC_PDB_METRIC tables contains records referenced to the pdb metric both with dbid of the CDB and dbid of the PDB
--(rows duplicated) will be used the dbid that owns more snapshot (it depends on the retention period, sometimes different from CDB and PDB and on awr_pdb_autoflush_enabled, 
--if not TRUE, snapshot with PDB dbid won't be present)
select dbid from 
(
	select dbid,count(*) from dba_hist_snapshot 
	where BEGIN_INTERVAL_TIME > trunc(sysdate-30) 
	group by dbid
	order by count(*) desc
)
where rownum = 1;

--On DBA_HIST_CON_SYSMETRIC_SUMM and DBA_HIST_RSRC_PDB_METRIC even if it is possible to use dbid referred to the CDB the CON_ID must be the one referred to the PDB
SELECT SYS_CONTEXT('USERENV','CON_ID') con_id FROM dual;

--Retrieve min and max snap_id to use considering only SNAPDAYS days
SELECT nvl(MIN(snap_id),0) bid, nvl(MAX(snap_id),0) eid, nvl(extract(day from MAX(END_INTERVAL_TIME)-MIN(BEGIN_INTERVAL_TIME)),0) snapdaysretrieved
FROM dba_hist_snapshot
WHERE BEGIN_INTERVAL_TIME > trunc(sysdate-&&SNAPDAYS) 
AND instance_number = &V_INSTANCE_NUMBER AND dbid=&V_DBID;

declare
current_day_difference number := 0;
--Monthly
snap_count_month number := 0;
day_month_count number := 0;
--"Old date" used to have a comparison value
reference_date date := trunc(to_date('01/01/2000', 'dd/mm/yyyy'));
current_date_changed_month date := reference_date;
current_date_month date;
v_cpu_db_month_avg number := 0;
v_cpu_db_month_max number := 0;
v_IOPS_month_avg number := 0;
v_IOMB_month_avg number := 0;
v_IOMB_month_max number := 0;
--Weekly
type arrayOfNumber_weeks is VARRAY(4) of number;
type arrayOfDates_weeks is VARRAY(4) of date;
current_date_week date;
--Positional array for the last 4 weeks
current_date_changed_week arrayOfDates_weeks := arrayOfDates_weeks(reference_date,reference_date,reference_date,reference_date);
day_week_count arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
snap_count_week arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
v_cpu_db_weekly_avg arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
v_cpu_db_weekly_max arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
v_IOPS_weekly_avg arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
v_IOMB_weekly_avg arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
v_IOMB_weekly_max arrayOfNumber_weeks := arrayOfNumber_weeks(0,0,0,0);
--Weekly averages
type arrayOfNumber_days is VARRAY(7) of number;
current_week number := -1;
--Positional array for each day of the last week
daily_count arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
v_cpu_db_daily_avg arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
v_cpu_db_daily_max arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
v_IOPS_daily_avg arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
v_IOMB_daily_avg arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
v_IOMB_daily_max arrayOfNumber_days := arrayOfNumber_days(0,0,0,0,0,0,0);
--Daily series
time_series_count number := 0;
TYPE arrayOfClob_table IS TABLE OF clob;
v_time_series_table arrayOfClob_table := arrayOfClob_table();
v_time_series_element clob := '';

begin
	--if the DIAGNOSTIC PACK HAS ALREADY BEEN USED AND THE DB HAS AT LEAST 1 DAY OF SNAPSHOTS
	if (&V_diagnostic_pack_utilizations>0 AND &V_SNAP_DAYS_RETRIEVED>0) then
		for metric_row in (
			select snap_id, max(end_time) end_time,
			nvl(max(decode(metric_name, 'CPU Usage Per Sec', round(average/100,2),null)), 0) cpu_per_s,
			nvl(max(decode(metric_name, 'CPU Usage Per Sec', round(maxval/100,2),null)), 0) cpu_per_s_max,
			--Not valued for PDB
			--nvl(max(decode(metric_name, 'Host CPU Usage Per Sec', round(average/100,2),null)), 0) h_cpu_per_s,
			--nvl(max(decode(metric_name, 'Host CPU Usage Per Sec', round(maxval/100,2),null)), 0) h_cpu_per_s_max,
			--nvl(max(decode(metric_name, 'Physical Read Total IO Requests Per Sec', average,null)),0) read_iops,
			--nvl(max(decode(metric_name, 'Physical Read Total IO Requests Per Sec', maxval,null)),0) read_iops_max,
			nvl(max(decode(metric_name, 'Physical Read Total Bytes Per Sec', round((average)/1024/1024,2),null)),0) read_mb_s,
			nvl(max(decode(metric_name, 'Physical Read Total Bytes Per Sec', round((maxval)/1024/1024,2),null)),0) read_mb_s_max,
			--Not valued for PDB
			--nvl(max(decode(metric_name, 'Physical Write Total IO Requests Per Sec', average,null)),0) write_iops,
			--nvl(max(decode(metric_name, 'Physical Write Total IO Requests Per Sec', maxval,null)),0) write_iops_max,
			nvl(max(decode(metric_name, 'Physical Write Total Bytes Per Sec', round((average)/1024/1024,2),null)),0) write_mb_s,
			nvl(max(decode(metric_name, 'Physical Write Total Bytes Per Sec', round((maxval)/1024/1024,2),null)),0) write_mb_s_max,
			max(decode(metric_name,'PDB_IOPS', average,null)) pdb_iops
			--Not used, decided to use read_mb_s+write_mb_s and cpu_per_s
			--max(decode(metric_name,'PDB_IOMBPS', average,null)) "PDB_IOMBPS",
			--max(decode(metric_name,'PDB_CPU_USAGE_PER_S', average,null)) "PDB_CPU_USAGE_PER_S"	
			from(
					select snap_id, end_time, metric_name, average, maxval
					from DBA_HIST_CON_SYSMETRIC_SUMM
					where dbid = &V_DBID and con_id = &V_CONID and snap_id between &V_BID and &V_EID and instance_number=&V_INSTANCE_NUMBER 
					and	metric_name in ('CPU Usage Per Sec','Host CPU Usage Per Sec','Physical Read Total Bytes Per Sec',
					'Physical Read Total IO Requests Per Sec','Physical Write Total Bytes Per Sec','Physical Write Total IO Requests Per Sec')
					--Part dedicated to PDB 
					union all (
						select snap_id, cast(end_time as DATE), metric_name, round(metric_value,1) average, null maxval
						from
						(  
							select snap_id, dbid, con_id, end_time,instance_number, 'PDB_IOPS' metric_name, IOPS metric_value from DBA_HIST_RSRC_PDB_METRIC 
							--union all
							--select snap_id, dbid, end_time,instance_number, 'PDB_IOMBPS' metric_name, IOMBPS metric_value from DBA_HIST_RSRC_PDB_METRIC 
							--union all
							--select snap_id, dbid, end_time,instance_number, 'PDB_CPU_USAGE_PER_S' metric_name, CPU_CONSUMED_TIME/INTSIZE_CSEC/10 metric_value from DBA_HIST_RSRC_PDB_METRIC 							
						)
						where dbid = &V_DBID and con_id = &V_CONID and snap_id between &V_BID and &V_EID and instance_number=&V_INSTANCE_NUMBER						
						and snap_id -1 not in 
							(
								select max(snap_id) last_snap_before_seq_chg from DBA_HIST_RSRC_PDB_METRIC
								where dbid = &V_DBID and con_id = &V_CONID and snap_id between &V_BID and &V_EID and instance_number=&V_INSTANCE_NUMBER
								group by sequence#
								--PDB_IOPS and PDB_IOMBPS are wrong when sequence changes in DBA_HIST_RSRC_PDB_METRIC
								--remove those snapshot right after a change in sequence (this is the "snap_id -1 not in" subquery)
							)
					)				
			)
			group by snap_id
			order by snap_id
		)
		
		loop	
			current_day_difference := trunc(sysdate)-trunc(metric_row.end_time);
			--Calculating monthly averages and maximums
			if current_day_difference < &&MONTH then
				--Counter to know how many effective days have snapshots (it could be the case that in 30 days interval, snapshots are not available for x days?)
				current_date_month := trunc(metric_row.end_time);
				if current_date_changed_month != current_date_month then
					current_date_changed_month := current_date_month;
					day_month_count := day_month_count+1;
				end if;
				--number of snapshot considered during last month
				snap_count_month := snap_count_month + 1;
				v_cpu_db_month_avg := v_cpu_db_month_avg + nvl(round(metric_row.cpu_per_s,2),0);		
				v_IOPS_month_avg := v_IOPS_month_avg + nvl(round((metric_row.pdb_iops),2),0);		
				v_IOMB_month_avg := v_IOMB_month_avg + nvl(round((metric_row.read_mb_s+metric_row.write_mb_s),2),0);
				if(metric_row.cpu_per_s_max>=v_cpu_db_month_max) then 
					v_cpu_db_month_max := round(metric_row.cpu_per_s_max,2); 
				end if;
				if((metric_row.read_mb_s_max+metric_row.write_mb_s_max)>=v_IOMB_month_max) then
					v_IOMB_month_max := round((metric_row.read_mb_s_max+metric_row.write_mb_s_max),2);
				end if;
			end if;
			
			--Integer that represents the week number (0 current, 1 next week, ..., 3)
			current_week := trunc(current_day_difference/&&WEEK);
			--Consider only last 4 weeks
			if current_week < &&WEEKS_NUMBER then
				--Counter to know how many effective days have snapshots (it could be the case that in 7 days interval, snapshots are not available for x days?)				
				current_date_week := trunc(metric_row.end_time);
				if current_date_changed_week(current_week+1) != current_date_week then
					current_date_changed_week(current_week+1) := current_date_week;
					day_week_count(current_week+1) := day_week_count(current_week+1)+1;
				end if;	
				--number of snapshot considered during the current week		
				snap_count_week(current_week+1) := snap_count_week(current_week+1) + 1;				
				v_cpu_db_weekly_avg(current_week+1) := v_cpu_db_weekly_avg(current_week+1) + nvl(round(metric_row.cpu_per_s,2),0);
				v_IOPS_weekly_avg(current_week+1) := v_IOPS_weekly_avg(current_week+1) + nvl(round((metric_row.pdb_iops),2),0);
				v_IOMB_weekly_avg(current_week+1) := v_IOMB_weekly_avg(current_week+1) + nvl(round(metric_row.read_mb_s+metric_row.write_mb_s,2),0);		
				if(metric_row.cpu_per_s_max>=v_cpu_db_weekly_max(current_week+1)) then 
					v_cpu_db_weekly_max(current_week+1) := round(metric_row.cpu_per_s_max,2); 
				end if;				
				if((metric_row.read_mb_s_max+metric_row.write_mb_s_max)>=v_IOMB_weekly_max(current_week+1)) then
					v_IOMB_weekly_max(current_week+1) := round((metric_row.read_mb_s_max+metric_row.write_mb_s_max),2);
				end if;			
			end if;
			
			--Calculating daily averages and highs for the last 7 days
			--If the snapshot is from the last 7 days (today->day_number=0)
			if current_day_difference < &&WEEK then		
				if current_day_difference<7 then 
					--the array starts from position 1 (today->position 1 in the array, 7 days ago->position 7 in the array)
					daily_count(current_day_difference+1) := daily_count(current_day_difference+1) + 1;
					v_cpu_db_daily_avg(current_day_difference+1) := v_cpu_db_daily_avg(current_day_difference+1) + nvl(round(metric_row.cpu_per_s,2),0);					
					v_IOPS_daily_avg(current_day_difference+1) := v_IOPS_daily_avg(current_day_difference+1) + nvl(round((metric_row.pdb_iops),2),0);
					v_IOMB_daily_avg(current_day_difference+1) := v_IOMB_daily_avg(current_day_difference+1) + nvl(round(metric_row.read_mb_s+metric_row.write_mb_s,2),0);					
					if(metric_row.cpu_per_s_max>=v_cpu_db_daily_max(current_day_difference+1)) then
						v_cpu_db_daily_max(current_day_difference+1) := round(metric_row.cpu_per_s_max,2); 
					end if;
					if(metric_row.read_mb_s_max+metric_row.write_mb_s_max>=v_IOMB_daily_max(current_day_difference+1)) then
						v_IOMB_daily_max(current_day_difference+1) := round(metric_row.read_mb_s_max+metric_row.write_mb_s_max,2);
					end if;
				end if;
			end if;
			
			--Daily data by series
			if sysdate - metric_row.end_time < 1 then
				time_series_count := time_series_count + 1;
				v_time_series_element := '';
				v_time_series_element := concat(v_time_series_element, concat(to_char(metric_row.end_time, 'ddmmHH24:MI'), '|||'));
				v_time_series_element := concat(v_time_series_element, concat(nvl(to_char(round(metric_row.cpu_per_s,2)),'N/A'), '|||'));
				v_time_series_element := concat(v_time_series_element, concat(nvl(to_char(round(metric_row.cpu_per_s_max,2)),'N/A'), '|||'));	
				v_time_series_element := concat(v_time_series_element, concat(nvl(to_char(round(metric_row.pdb_iops,2)),'N/A'), '|||'));
				v_time_series_element := concat(v_time_series_element, concat(nvl(to_char(round(metric_row.read_mb_s+metric_row.write_mb_s,2)),'N/A'), '|||'));
				v_time_series_element := concat(v_time_series_element, nvl(to_char(round(metric_row.read_mb_s_max+metric_row.write_mb_s_max,2)),'N/A'));
				v_time_series_table.EXTEND;
				v_time_series_table(v_time_series_table.LAST) := v_time_series_element;
			end if;	   
		end loop;
		
		--PLACEHOLDER begin important output
		DBMS_OUTPUT.PUT_LINE('BEGINOUTPUT');

		--Monthly output (show only if the number of days for which you have snapshots is greater than MONTH_SOFT_LIMIT [arbitrarily decided soft limit])
		if(day_month_count>=&&MONTH_SOFT_LIMIT) then 
			--Avg calculated on monthly snapshot number
			if(snap_count_month>0) then
				v_cpu_db_month_avg := round((v_cpu_db_month_avg / snap_count_month),2);
				v_IOPS_month_avg := round((v_IOPS_month_avg / snap_count_month),2);
				v_IOMB_month_avg := round((v_IOMB_month_avg / snap_count_month),2);
			end if;
			--v_cpu_db_month_avg,v_cpu_db_month_max,v_IOPS_month_avg,v_IOMB_month_avg,v_IOMB_month_max
			DBMS_OUTPUT.PUT_LINE(v_cpu_db_month_avg || '|||' || v_cpu_db_month_max || '|||' || v_IOPS_month_avg || '|||' || v_IOMB_month_avg || '|||' || v_IOMB_month_max);
		else
			DBMS_OUTPUT.PUT_LINE('N/A|||N/A|||N/A|||N/A|||N/A');
		end if;
		
		--Output last 4 weeks averages and maximums
		--Loops over weeks
		for w_index in day_week_count.FIRST..day_week_count.LAST
		LOOP
			--Weekly output (show only if the number of days for which you have snapshots is greater than WEEK_SOFT_LIMIT [arbitrarily decided soft limit])
			if (day_week_count(w_index)>=&&WEEK_SOFT_LIMIT) THEN
				if  (snap_count_week(w_index)>0) THEN
					v_cpu_db_weekly_avg(w_index) := round((v_cpu_db_weekly_avg(w_index) / snap_count_week(w_index)),2);
					v_IOPS_weekly_avg(w_index) := round((v_IOPS_weekly_avg(w_index) / snap_count_week(w_index)),2);
					v_IOMB_weekly_avg(w_index) := round((v_IOMB_weekly_avg(w_index) / snap_count_week(w_index)),2);
					--v_cpu_db_weekly_avg(w_index),v_cpu_db_weekly_max(w_index),v_IOPS_weekly_avg(w_index),v_IOMB_weekly_avg(w_index),v_IOMB_weekly_max(w_index)
					DBMS_OUTPUT.PUT_LINE(v_cpu_db_weekly_avg(w_index) || '|||' || v_cpu_db_weekly_max(w_index) || '|||' || v_IOPS_weekly_avg(w_index) || '|||' || v_IOMB_weekly_avg(w_index) || '|||' || v_IOMB_weekly_max(w_index));
				ELSE
					DBMS_OUTPUT.PUT_LINE('N/A|||N/A|||N/A|||N/A|||N/A');
				end if;
			ELSE
				DBMS_OUTPUT.PUT_LINE('N/A|||N/A|||N/A|||N/A|||N/A');
			end if;
		END LOOP;
	
		--Output daily averages last week
		--Daily averages, loops over days
		for l_index in daily_count .FIRST..daily_count.LAST
		loop 
		--the array starts from position 1 (today->position 1 in the array, 7 days ago->position 7 in the array)
			if(daily_count(l_index)>0) then
				v_cpu_db_daily_avg(l_index) := round(v_cpu_db_daily_avg(l_index) / daily_count(l_index),2);
				v_IOPS_daily_avg(l_index) := round(v_IOPS_daily_avg(l_index) / daily_count(l_index),2);
				v_IOMB_daily_avg(l_index) := round(v_IOMB_daily_avg(l_index) / daily_count(l_index),2);
				--v_cpu_db_daily_avg(l_index),v_cpu_db_daily_max(l_index),v_IOPS_daily_avg(l_index),v_IOMB_daily_avg(l_index),v_IOMB_daily_max(l_index)
				DBMS_OUTPUT.PUT_LINE(v_cpu_db_daily_avg(l_index) || '|||' || v_cpu_db_daily_max(l_index) || '|||' || v_IOPS_daily_avg(l_index) || '|||' || v_IOMB_daily_avg(l_index) || '|||' || v_IOMB_daily_max(l_index));
			else
				DBMS_OUTPUT.PUT_LINE('N/A|||N/A|||N/A|||N/A|||N/A');
			end if;		
		END LOOP;
		
		--Daily series
		FOR v_time_series_table_index IN v_time_series_table.FIRST..v_time_series_table.LAST
		LOOP
			DBMS_OUTPUT.PUT_LINE(v_time_series_table(v_time_series_table_index));
		END LOOP;	
		
		--PLACEHOLDER end important output
		DBMS_OUTPUT.PUT_LINE('ENDOUTPUT');
	--DIAGNOSTIC PACK never used
	else
		DBMS_OUTPUT.PUT_LINE('BEGINOUTPUT');
		DBMS_OUTPUT.PUT_LINE('ENDOUTPUT');
	end if;
end;
/

exit