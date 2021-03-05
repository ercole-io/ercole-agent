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
SELECT 
  a.SCHEMA_NAME as database_name,
  a.DEFAULT_CHARACTER_SET_NAME as charset,
  a.DEFAULT_COLLATION_NAME as collation,
  a.DEFAULT_ENCRYPTION as encrypted
FROM 
  information_schema.schemata AS a;
