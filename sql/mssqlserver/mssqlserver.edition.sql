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
	AS 'edition'