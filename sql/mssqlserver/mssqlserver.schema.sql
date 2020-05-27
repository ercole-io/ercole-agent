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
;WITH [ONE] AS (
	SELECT
		/*ind.type, ---0 = Heap, 1 = Clustered, > 1 Non clustred indexes
		ind.type_desc,
		au.type_desc, --0 = Dropped, 1 = In-row data (all data types, except LOB data types), 2 = Large object (LOB) data (text, ntext, image, xml, large value types, and CLR user-defined types), 3 = Row-overflow data
		au.type,*/
		CASE 
			WHEN (ind.type in (0) and au.type in (1, 3)) THEN 'Table Data'
			WHEN (ind.type in (1) and au.type in (1, 3)) THEN 'Clustered Index Data'
			WHEN (ind.type >1 and au.type in (1, 3)) THEN 'Non Clustered Index Data'
			ELSE 'LOB data' 
		END as [allocation_type],
		au.used_pages,
		au.total_pages
	FROM
		sys.objects obj
		inner join sys.indexes ind 
			on obj.object_id = ind.object_id
		inner join sys.partitions part 
		    on ind.object_id = part.object_id and ind.index_id = part.index_id
		inner join sys.allocation_units au
		    on part.partition_id = au.container_id
)
SELECT
	@@SERVERNAME AS [servername],
	DB_ID() AS [database_id],
	DB_NAME() AS [database_name],
	[ONE].[allocation_type],
	convert( decimal(36,3), SUM( ([ONE].[used_pages] * (8E/1024)) ) ) AS [used_mb],
	convert( decimal(36,3), SUM( ([ONE].[total_pages] * (8E/1024)) ) ) AS [allocated_mb]
FROM
	[ONE]
GROUP BY 
	[ONE].[allocation_type]