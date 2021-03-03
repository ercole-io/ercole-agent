#-- Copyright (c) 2019 Sorint.lab S.p.A.
#
#-- This program is free software: you can redistribute it and/or modify
#-- it under the terms of the GNU General Public License as published by
#-- the Free Software Foundation, either version 3 of the License, or
#-- (at your option) any later version.
#
#-- This program is distributed in the hope that it will be useful,
#-- but WITHOUT ANY WARRANTY; without even the implied warranty of
#-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#-- GNU General Public License for more details.
#
#-- You should have received a copy of the GNU General Public License
#-- along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
WITH ONE AS(
  SELECT
    TABLE_SCHEMA,
	ENGINE,
    round(sum(DATA_LENGTH+INDEX_LENGTH)/1024/1024,3) as ALLOC_MB
  FROM 
    information_schema.tables 
  WHERE     
    ENGINE='MyISAM'    
  GROUP BY     
   TABLE_SCHEMA
), TWO as(
  SELECT 
    SUBSTRING_INDEX(tbs.NAME,'/',1) as TABLE_SCHEMA,
    round(sum(tbs.allocated_size)/1024/1024,3) as ALLOC_MB
  FROM     
    information_schema.INNODB_TABLESPACES as tbs
  WHERE      
    FILE_SIZE>0    
  GROUP BY 
    TABLE_SCHEMA
)
SELECT
  concat(@@global.hostname,':',@@global.port) as instance,
  TWO.TABLE_SCHEMA,
  'InnoDB' AS ENGINE,
  TWO.ALLOC_MB
FROM
  TWO
UNION ALL
SELECT 
  concat(@@global.hostname,':',@@global.port) as instance,
  TABLE_SCHEMA,
  ENGINE,
  ALLOC_MB
FROM 
  ONE
ORDER BY
  TABLE_SCHEMA,
  ENGINE;
