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

VARIABLE STATUS varchar2(100);
VARIABLE PSU_DATE varchar2(100);
VARIABLE DESCRIPTION varchar2(100);
VARIABLE VERSION varchar2(100);
VARIABLE EXTENDVERSION varchar2(100);
VARIABLE EXIST number;

set lines 8000 pages 0 feedback off verify off
set colsep "|||"
alter session set NLS_DATE_FORMAT='YYYY-MM-DD';

BEGIN
select count(*) into :EXIST  from registry$sqlpatch;
SELECT DBMS_DB_VERSION.VERSION || '.' || DBMS_DB_VERSION.RELEASE into :VERSION FROM v$instance;
-- 12.1
   IF ( :VERSION = '12.1' AND :EXIST > 0 ) THEN 
   	  	with PSU as
		(
			select DESCRIPTION as DESCRIPTION
			from  registry$sqlpatch
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and upper(DESCRIPTION) like 'DATABASE BUNDLE PATCH%'
					or upper(DESCRIPTION) like 'DATABASE PATCH SET UPDATE%'
				   )
		),
		DATA as
		(
			select 
			case
		    WHEN length(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1)) >= 6 THEN TO_DATE(substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,7),'YYMMDD')
			WHEN substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,1) = '1' THEN TO_DATE('141014','YYMMDD') 
			WHEN substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,1) = '2' THEN TO_DATE('150120','YYMMDD') 
			WHEN substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,1) = '3' THEN TO_DATE('150414','YYMMDD') 
			WHEN substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,1) = '4' THEN TO_DATE('150714','YYMMDD') 
			WHEN substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,1) = '5' THEN TO_DATE('151020','YYMMDD') 
				ELSE
					TO_DATE(substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,7),'YYMMDD')
			END as PSU_DATE
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
								  from  registry$sqlpatch 
								  where ACTION='APPLY' 
								  and upper(DESCRIPTION) like 'DATABASE BUNDLE PATCH%'
								  or upper(DESCRIPTION) like 'DATABASE PATCH SET UPDATE%'
	  							 )
		),
		STATE as
		(
			select 
			case   
				WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
				ELSE 'OK' 
			END as "STATUS"
			from DATA
		)
		select DESCRIPTION,PSU_DATE,STATUS 
		into :DESCRIPTION,:PSU_DATE,:STATUS 
		from PSU,DATA,STATE;
-- 12.2   
   ELSIF ( :VERSION = '12.2' AND :EXIST > 0 ) THEN 
   	  	with PSU as
		(
			select DESCRIPTION as DESCRIPTION
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and upper(DESCRIPTION) like 'DATABASE%RELEASE UPDATE%'
				   )
		),
		DATA as
		(
			select TO_DATE(substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),1,7),'YYMMDD') as PSU_DATE
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and upper(DESCRIPTION) like 'DATABASE%RELEASE UPDATE%'
				   )
		),
		STATE as
		(
			select 
			case   
				WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
				ELSE 'OK' 
			END as "STATUS"
			from DATA
		)
		select DESCRIPTION,PSU_DATE,STATUS 
		into :DESCRIPTION,:PSU_DATE,:STATUS 
		from PSU,DATA,STATE;
   ELSIF ( :VERSION = '18.0' AND :EXIST > 0) THEN 
     	with PSU as
		(
			select DESCRIPTION as DESCRIPTION
			from  registry$sqlpatch
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and DESCRIPTION like '%Release Update%'
				   )
		),
		DATA as
		(
			select TO_DATE(Substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1), 1,instr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),' ') - 1),'YYMMDD') as PSU_DATE
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and DESCRIPTION like '%Release Update%'
				   )
		),
		STATE as
		(
			select 
			case   
				WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
				ELSE 'OK' 
			END as "STATUS"
			from DATA
		)
		select DESCRIPTION,PSU_DATE,STATUS 
		into :DESCRIPTION,:PSU_DATE,:STATUS 
		from PSU,DATA,STATE;
   ELSIF ( :VERSION = '19.0' AND :EXIST > 0) THEN 
     	with PSU as
		(
			select DESCRIPTION as DESCRIPTION
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and DESCRIPTION like '%Release Update%'
				   )
		),
		DATA as
		(
			select TO_DATE(Substr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1), 1,instr(substr(DESCRIPTION, - instr(reverse(DESCRIPTION), '.') + 1),' ') - 1),'YYMMDD') as PSU_DATE
			from  registry$sqlpatch 
			where action_time = ( select max(action_time) 
					from  registry$sqlpatch 
					where ACTION='APPLY' 
					and DESCRIPTION like '%Release Update%'
				   )
		),
		STATE as
		(
			select 
			case   
				WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
				ELSE 'OK' 
			END as "STATUS"
			from DATA
		)
		select DESCRIPTION,PSU_DATE,STATUS 
		into :DESCRIPTION,:PSU_DATE,:STATUS 
		from PSU,DATA,STATE;
    ELSE
     	select 'N/A','N/A','N/A' 
		into :DESCRIPTION,:PSU_DATE,:STATUS 
		from dual;
   END IF; 

END;
/

col Description for a70
col PSU for a40
col STATUS for a40
select :DESCRIPTION as Description 
	   ,:PSU_DATE as PSU
--	   ,:STATUS as STATUS
from dual WHERE :PSU_DATE != 'N/A';