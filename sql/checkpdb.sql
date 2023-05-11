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

SET lines 32767 pages 0 feedback OFF verify OFF timing off
SET colsep "|||"

SELECT CASE
           WHEN COUNT(*) > 0 THEN 'TRUE'
           WHEN count(*) = 0 THEN 'FALSE'
       END
FROM v$pdbs where NAME not like 'PDB$SEED' ;
exit