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
	DB_ID() AS [database_id],
	DB_NAME() AS [database_name],
	[df].file_id as [file_id],
	[df].name as [file_name], 
	round( ([df].size*(8E/1024)),2) as [allocated_mb],
	round( ([df].size*(8E/1024)) - ( fileproperty([df].name,'SpaceUsed')*(8E/1024) ) ,3) AS [freeSpace_mb],
	substring([df].type_desc,1,1)  as [fileType],
	[df].state_desc as [status]	
FROM 
	sys.database_files as [df]

