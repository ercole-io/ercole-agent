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
    TABLE_NAME, 
    ENGINE,
    round(sum(DATA_LENGTH+INDEX_LENGTH)/1024/1024,3) as ALLOC_MB,
    round(sum(DATA_LENGTH)/1024/1024,3) as DATA_MB, 
    round(sum(INDEX_LENGTH)/1024/1024,3) as INDEX_MB,
    round(sum(DATA_FREE)/1024/1024,3) as FREE_MB
  FROM 
    information_schema.tables 
  WHERE     
    ENGINE='MyISAM'    
  GROUP BY     
   TABLE_SCHEMA, 
   TABLE_NAME, 
   ENGINE
), TWO as(
  SELECT 
  SUBSTRING_INDEX(tbs.NAME,'/',1) as TABLE_SCHEMA,
  SUBSTRING_INDEX(tbs.NAME,'/',-1) as TABLE_NAME,
  tbs.PAGE_SIZE,
  round(sum(tbs.allocated_size)/1024/1024,3) as ALLOC_MB
  FROM     
  information_schema.INNODB_TABLESPACES as tbs
  WHERE      
  FILE_SIZE>0    
  GROUP BY 
  `TABLE_SCHEMA`, `TABLE_NAME`, `PAGE_SIZE`
), THREE AS (
  SELECT
    TWO.TABLE_SCHEMA,
    TWO.TABLE_NAME,
    'InnoDB' as ENGINE,
	TWO.ALLOC_MB,
    ( SELECT 
        IFNULL(round((sum(stat_value) * TWO.PAGE_SIZE)/1024/1024,3),0)
      FROM 
        mysql.innodb_index_stats 
      WHERE 
        stat_name='size' 
        AND database_name=TWO.TABLE_SCHEMA
        AND table_name=TWO.TABLE_NAME
    ) AS INDEX_MB,
    ( SELECT
        round(sum(DATA_FREE)/1024/1024,3) as FREE_MB
      FROM
        information_schema.tables
      WHERE 
        ENGINE = 'InnoDB'
        AND TABLE_SCHEMA=TWO.TABLE_SCHEMA
        AND TABLE_NAME=TWO.TABLE_NAME
      GROUP BY
        TABLE_SCHEMA,
        TABLE_NAME
      ) as FREE_MB
  FROM
    TWO
)
SELECT
    THREE.TABLE_SCHEMA,
    THREE.TABLE_NAME,
    'InnoDB' as ENGINE,
    THREE.ALLOC_MB,
    (THREE.ALLOC_MB - THREE.INDEX_MB) AS DATA_MB,
    IFNULL(THREE.INDEX_MB, 0) AS INDEX_MB,
    IFNULL(THREE.FREE_MB, 0) AS FREE_MB
FROM
  THREE
UNION ALL
SELECT 
  TABLE_SCHEMA,
  TABLE_NAME, 
  ENGINE,
  ALLOC_MB,  
  DATA_MB, 
  INDEX_MB,
  FREE_MB
FROM 
  ONE
ORDER BY 
  TABLE_SCHEMA,
  TABLE_NAME;
