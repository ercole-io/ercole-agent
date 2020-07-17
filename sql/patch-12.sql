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

column "Nome_Acronimo" for a8
column "DB_Name" for a10

select 
	   (select host_name from v$instance) as Hostname,
	   'ND',
	   (select value from v$parameter where name='db_name') as Nome_DB,
	   (select db_unique_name from v$database) as DB_Unique_name,
           ( select version from v$instance) as Version,
	   PATCH_ID,ACTION,DESCRIPTION,to_char(action_time,'YYYY-MM-DD') from  registry$sqlpatch order by action_time;
exit
