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
		database_name,
		Rounded = CAST(DATEADD(MINUTE, ROUND(DATEDIFF(MINUTE, 0, [backup_start_date]) / 15.0, 0) * 15, 0) as TIME(0)),
		[type],
		[start_wd] = datepart(WEEKDAY, backup_start_date),
		[backup_size_gb] = (backup_size/1024/1024/1024)
	FROM
		msdb.dbo.backupset as bkp
		JOIN master.sys.databases as db on (db.name = bkp.database_name)
	WHERE
		db.database_id = @pdatabase_id
		and bkp.backup_start_date >= cast(Dateadd(month,-1,getdate()) as date)
		and bkp.is_copy_only = 0
)
, [TWO] AS (
SELECT 
	[ONE].database_name, 
	CASE [ONE].type 
		WHEN 'D' THEN 'Database'
		WHEN 'I' THEN 'Differential database'
		WHEN 'L' THEN 'Log'
		WHEN 'F' THEN 'File or filegroup'
		WHEN 'G' THEN 'Differential file'
		WHEN 'P' THEN 'Partial'
		WHEN 'Q' THEN 'Differential partial'
	END as [type],
	[ONE].Rounded,
	CASE [ONE].start_wd
			WHEN 1 THEN 'Sunday'
			WHEN 2 THEN 'Monday'
			WHEN 3 THEN 'Tuesday'
			WHEN 4 THEN 'Wednesday'
			WHEN 5 THEN 'Thursday'
			WHEN 6 THEN 'Friday'
			WHEN 7 THEN 'Saturday'
		END as start_wd,
	ROUND(AVG([ONE].backup_size_gb),3) AS backup_size_gb,
	count(*) as qtd
FROM 
	[ONE]
GROUP BY 
	[ONE].database_name, 
	[ONE].type, 
	[ONE].start_wd, 
	[ONE].Rounded
) 
SELECT [TWO].database_name, 
	[TWO].type as [backup_type],  
	[TWO].Rounded as [hour],
	CAST(ROUND(AVG([TWO].backup_size_gb),3)  AS DECIMAL(12,3)) as [avg_bck_size_gb],
	--STRING_AGG([TWO].start_wd, ',') as [DaysOfWeek] --only for Version >= 2016
	SUBSTRING(
        (
            SELECT 
				','+ CHILD.start_wd  AS [text()]
            FROM 
				[TWO] AS [CHILD]
            WHERE 
				[CHILD].database_name = [TWO].database_name
				AND [CHILD].type = [TWO].type
				AND [CHILD].Rounded = [TWO].Rounded
            ORDER BY 
				[CHILD].database_name,
				[CHILD].type,
				[CHILD].Rounded
            FOR XML PATH ('')
        ), 2, 1000) [week_days]
FROM
	[TWO]
GROUP BY 
	[TWO].database_name, 
	[TWO].type,  
	[TWO].Rounded
ORDER BY
	[TWO].database_name, 
	[TWO].type,  
	[TWO].Rounded