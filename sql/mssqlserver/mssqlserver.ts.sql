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
	SERVERPROPERTY('ComputerNamePhysicalNetBIOS') AS [hostname],
	@@SERVERNAME AS [servername],
	DB_ID() AS [database_id],
	DB_NAME() AS [database_name],
	[df].NAME as [file_name], 
	round( ([df].size*(8E/1024)),2) as [alloc_mb],
	round( (fileproperty([df].name,'SpaceUsed')*(8E/1024) ),2)  AS [used_mb],
		round( ( fileproperty([df].name,'SpaceUsed')*(8E/1024) ) / ([df].size*(8E/1024)) ,2) *100 AS [used_percentage],
	case [df].[IS_PERCENT_GROWTH] when 1 then [df].[GROWTH] else ([df].[GROWTH] *(8E/1024)) end AS [growth], 
	case [df].is_percent_growth when 1 then '%' else 'MB' end as [growthUnit],
	[df].type_desc  as [fileType],
	[df].state_desc as [status]	
FROM sys.database_files as [df]

