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

set lines 8000 pages 0 feedback off verify off timing off
set colsep "|||"

define TABLENAME = 'dba_services'
col :TABLENAME_ new_val TABLENAME noprint
variable TABLENAME_ varchar2(30) 

define NAMECOLUMNS = '-1 pdb,name,FAILOVER_METHOD,FAILOVER_TYPE,FAILOVER_RETRIES,FAILOVER_DELAY,ENABLED'
col :NAMECOLUMNS_ new_val NAMECOLUMNS noprint
variable NAMECOLUMNS_ varchar2(100)

DECLARE 
DB_VERSION number;
begin
	DB_VERSION := dbms_db_version.version + (dbms_db_version.release / 10);
	if DB_VERSION < 12.1 then
		:TABLENAME_ := 'dba_services';
		:NAMECOLUMNS_ := ''''' pdb,name,FAILOVER_METHOD,FAILOVER_TYPE,FAILOVER_RETRIES,FAILOVER_DELAY,ENABLED';
	else
		:TABLENAME_ := 'cdb_services';
		:NAMECOLUMNS_ := 'pdb,name,FAILOVER_METHOD,FAILOVER_TYPE,FAILOVER_RETRIES,FAILOVER_DELAY,ENABLED';
	end if;	
end;
/

select :TABLENAME_, :NAMECOLUMNS_ from dual;

select &NAMECOLUMNS from &TABLENAME where NAME not in('SYS$BACKGROUND','SYS$USERS');

exit
