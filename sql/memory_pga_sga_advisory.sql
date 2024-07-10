-- Copyright (c) 2024 Sorint.lab S.p.A.

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

set serveroutput on

DECLARE
has_memory_target_enabled number := 0;
memory_size_lower number := 0;
exists_low_pga_cachehit_100 number := 0;
exists_low_pga_cachehit_95 number := 0;
pga_target_lower number := 0;
exists_low_est_dbtime_105 number := 0;
sga_size_lower number := 0;
begin
--Check if exists a record in v$memory_target_advice -> if yes automatic memory management is enabled
select count(*) into has_memory_target_enabled from v$memory_target_advice;
if (has_memory_target_enabled>0) then
	select count(*) into exists_low_est_dbtime_105 from v$memory_target_advice where ESTD_DB_TIME_FACTOR<=1.05 and MEMORY_SIZE_FACTOR<1;
	if (exists_low_est_dbtime_105>0) then
		select MEMORY_SIZE into memory_size_lower from
		(select MEMORY_SIZE
		from v$memory_target_advice
		where ESTD_DB_TIME_FACTOR<=1.05 and MEMORY_SIZE_FACTOR<1
		order by ESTD_DB_TIME_FACTOR desc,MEMORY_SIZE_FACTOR asc)
		where rownum=1;
	end if;
	DBMS_OUTPUT.PUT_LINE('BEGINOUTPUT');
	memory_size_lower := round(memory_size_lower/1024,3);
	if (exists_low_est_dbtime_105=0) then
		DBMS_OUTPUT.PUT_LINE('MEMORY_SIZE_LOWER_GB|||' || 'N/A');
	else
		DBMS_OUTPUT.PUT_LINE('MEMORY_SIZE_LOWER_GB|||' || memory_size_lower);
	end if;
	DBMS_OUTPUT.PUT_LINE('ENDOUTPUT');
else
	--Check if exist a pga_aggregate_target size lower than the actual that have pga_cache_hit_precentage = 100 or at least >=95
	select count(*) into exists_low_pga_cachehit_100 from v$pga_target_advice where ESTD_PGA_CACHE_HIT_PERCENTAGE=100 and PGA_TARGET_FACTOR<1;
	select count(*) into exists_low_pga_cachehit_95 from v$pga_target_advice where ESTD_PGA_CACHE_HIT_PERCENTAGE>=95 and PGA_TARGET_FACTOR<1;
	if (exists_low_pga_cachehit_100>0) then
		select PGA_TARGET_FOR_ESTIMATE into pga_target_lower from
		(select PGA_TARGET_FOR_ESTIMATE
		from v$pga_target_advice 
		where ESTD_PGA_CACHE_HIT_PERCENTAGE=100 and PGA_TARGET_FACTOR<1
		order by ESTD_PGA_CACHE_HIT_PERCENTAGE,PGA_TARGET_FACTOR)
		where rownum=1;
	elsif (exists_low_pga_cachehit_95>0) then
		select PGA_TARGET_FOR_ESTIMATE into pga_target_lower from
		(select PGA_TARGET_FOR_ESTIMATE 
		from v$pga_target_advice 
		where ESTD_PGA_CACHE_HIT_PERCENTAGE>=95 and PGA_TARGET_FACTOR<1
		order by ESTD_PGA_CACHE_HIT_PERCENTAGE,PGA_TARGET_FACTOR)
		where rownum=1;
	end if;
	--Check if exist a sga_target size lower than the actual that have est_db_time_factor at least 0.05 % greater than the actual 
	select count(*) into exists_low_est_dbtime_105 from v$sga_target_advice where ESTD_DB_TIME_FACTOR<=1.05 and SGA_SIZE_FACTOR<1;
	if (exists_low_est_dbtime_105>0) then
		select SGA_SIZE into sga_size_lower from
		(select SGA_SIZE
		from v$sga_target_advice 
		where ESTD_DB_TIME_FACTOR<=1.05 and SGA_SIZE_FACTOR<1
		order by ESTD_DB_TIME_FACTOR desc,SGA_SIZE_FACTOR asc)
		where rownum=1;
	end if;
	DBMS_OUTPUT.PUT_LINE('BEGINOUTPUT');
	pga_target_lower := round(pga_target_lower/1024/1024/1024,3);
	if (exists_low_pga_cachehit_100=0 and exists_low_pga_cachehit_95=0) then
		DBMS_OUTPUT.PUT_LINE('PGA_TARGET_AGGREGATE_LOWER_GB|||' || 'N/A');
	else
		DBMS_OUTPUT.PUT_LINE('PGA_TARGET_AGGREGATE_LOWER_GB|||' || pga_target_lower);
	end if;
	sga_size_lower := round(sga_size_lower/1024,3);
	if (exists_low_est_dbtime_105=0) then
		DBMS_OUTPUT.PUT_LINE('SGA_SIZE_LOWER_GB|||' || 'N/A');
	else
		DBMS_OUTPUT.PUT_LINE('SGA_SIZE_LOWER_GB|||' || sga_size_lower);
	end if;
	DBMS_OUTPUT.PUT_LINE('ENDOUTPUT');
end if;
end;
/

exit