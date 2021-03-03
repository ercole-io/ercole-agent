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
  concat(@@global.hostname,':',@@global.port) as instance,
  @@global.version as version,
  CASE 
    WHEN @@global.version_comment LIKE '%Community%' THEN  'COMMUNITY'
    ELSE 'ENTERPRISE'
  END as edition,
  @@global.version_compile_os as platform,
  @@global.version_compile_machine as architecture,
  @@global.default_storage_engine AS engine,
  (SELECT variable_value FROM performance_schema.global_status WHERE variable_name = 'innodb_redo_log_enabled') as redo_log_enabled,
  @@global.character_set_server as charset_server,
  @@global.character_set_system as charset_system,
  (@@global.innodb_page_size/1024) as page_size_KB,
  @@global.innodb_thread_concurrency as threads_concurrency,
  (@@global.innodb_buffer_pool_size/1024/1024) as buffer_pool_size_MB,
  (@@global.innodb_log_buffer_size/1024/1024) as log_buffer_size_MB,
  (@@global.innodb_sort_buffer_size /1024/1024) as sort_buffer_size_MB,
  @@global.read_only as read_only;
