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

set feedback off heading off pages 0 serverout on verify off lines 1234 timing off
set colsep "|||"

select name,round(total_mb/1024) as total_gb,round(usable_file_mb/1024) free_gb,round(usable_file_mb/total_mb*100,2) as free_percentage  from v$asm_diskgroup where name in (
select substr(value, 2, length(value)-1) from v$parameter where name='db_recovery_file_dest'
union
select distinct substr(name, 2, instr(name,'/')-2) from v$datafile where name like '+%'
union
select distinct substr(name, 2, instr(name,'/')-2) from v$tempfile where name like '+%');

exit