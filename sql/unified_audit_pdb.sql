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

-- Copyright (c) 2022 Sorint.lab S.p.A.

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

set feedback off heading off pages 0 serverout on verify off lines 1234 timing off
set lines 8000 pages 0 feedback off verify off timing off
set colsep "|||"

alter session set container=&1;

define IS12ORABOVE = '1=0'
col :IS12ORABOVE_ new_val IS12ORABOVE noprint
variable IS12ORABOVE_ varchar2(100) 

define FIELDNAME = '1=0'
col :FIELDNAME_ new_val FIELDNAME noprint
variable FIELDNAME_ varchar2(30) 

define TABLENAME = '1=0'
col :TABLENAME_ new_val TABLENAME noprint
variable TABLENAME_ varchar2(100) 

DECLARE 
DB_VERSION number;
begin
	DB_VERSION := dbms_db_version.version + (dbms_db_version.release / 10);
	if DB_VERSION >= 12.0 then
		:IS12ORABOVE_ := '(SELECT value FROM v$option WHERE parameter = ''Unified Auditing'')=''TRUE''';
		:FIELDNAME_ := 'POLICY_NAME';
		:TABLENAME_ := 'audit_unified_enabled_policies';
	ELSE
		:IS12ORABOVE_ := '1=0';
		:FIELDNAME_ := '''''';
		:TABLENAME_ := 'dual';
	end if;	
end;
/

select :IS12ORABOVE_, :FIELDNAME_, :TABLENAME_  from dual;

select &FIELDNAME from &TABLENAME where &IS12ORABOVE;

exit
