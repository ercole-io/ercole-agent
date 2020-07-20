-- Copyright (c) 2020 Sorint.lab S.p.A.

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


SELECT
	@@SERVERNAME AS [servername],
	database_id AS [database_id],
	name AS [database_name],
	(state_desc) AS [state_desc], 
	(	convert(varchar(128),SERVERPROPERTY('ProductVersion')) + ' '+ 
		CASE 
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Enterprise%' THEN 'ENT'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Express%' THEN 'EXP'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Standard%' THEN 'STD'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Business%' THEN 'BI'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Developer%' THEN 'DEV'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Web%' THEN 'WEB'
			WHEN convert(varchar(128), SERVERPROPERTY('Edition')) LIKE '%Azure%' THEN 'AZU'
			ELSE 'ANOTHER'
		END
	) AS [version],
	'Windows' AS [platform],
	recovery_model_desc as [recovery_model],
	collation_name AS [collation_name],
	8192 AS [blocksize],
	(	SELECT 
			count(*) 
		FROM
			sys.dm_os_schedulers 
		WHERE 
			status='VISIBLE ONLINE'
	) AS [schedulers_count],
	(	SELECT
			value_in_use 
		FROM
			sys.configurations 
		WHERE
			name='affinity mask'
	) as [affinity_mask],
	(	SELECT
			value_in_use 
		FROM
			sys.configurations 
		WHERE
			name='min server memory (MB)'
	) as [min_server_memory],
	(	SELECT
			value_in_use 
		FROM
			sys.configurations 
		WHERE
			name='max server memory (MB)'
	) as [max_server_memory],
	(	SELECT
			value_in_use 
		FROM
			sys.configurations 
		WHERE
			name='cost threshold for parallelism'
	) as [ctp],
	(	SELECT
			value_in_use 
		FROM
			sys.configurations 
		WHERE
			name='max degree of parallelism'
	) as [maxdop],
	(	SELECT
			round( (SUM(size*(8E))/1024/1024), 3)
		FROM 
			sys.database_files 
	) AS [Alloc]
FROM 
	sys.databases
WHERE
	database_id = @pdatabase_id;