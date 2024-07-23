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

set lines 8000 pages 0 feedback off verify off timing off
set colsep "|||"

alter session set container=&1;

SELECT
  (
    SELECT round(sum(bytes/1024/1024/1024))
    FROM dba_segments
   ) AS Alloc, 
   (
        (SELECT round(sum(bytes/1024/1024/1024))
        FROM dba_data_files)+
        (SELECT round(sum(bytes/1024/1024/1024))
        FROM dba_temp_files)+
        (SELECT round(sum(bytes/1024/1024/1024))
        FROM v$log)
    ),
    (
        (SELECT round(sum(decode(autoextensible,'NO',bytes/1024/1024/1024,'YES',maxbytes/1024/1024/1024)))
        FROM dba_data_files)+
        (SELECT round(sum(bytes/1024/1024/1024))
        FROM dba_temp_files)+
        (SELECT round(sum(bytes/1024/1024/1024))
        FROM v$log)
    ),
	(
		select value/1024/1024/1024 from v$system_parameter where name = 'sga_target'
	),
	(
		select value/1024/1024/1024 from v$system_parameter where name = 'pga_aggregate_target'
	)
FROM dual;

exit