--tested for PostgreSQL 9 and EPAS 10.
--all the data below refer to the db the client is connected to
WITH
--this query returns the size of the db
 a  AS (SELECT pg_database_size(current_database()) AS db_size)

--this query returns the number of all the non-system tables in the db
,d  AS (SELECT COUNT(*) AS tables_count FROM pg_tables WHERE schemaname NOT LIKE 'pg_%' AND
                                                             schemaname NOT LIKE 'dbms_%' AND
                                                             schemaname NOT LIKE 'aq$%' AND
                                                             schemaname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                             schemaname NOT LIKE 'utl_%' )

--this query returns the number of all the non-system indexes in the db
,e  AS (SELECT COUNT(*) AS indexes_count FROM pg_indexes WHERE schemaname NOT IN ('pg_catalog','sys'))

--this query returns the size of all the non-system tables in the db
,f  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_tables WHERE schemaname NOT LIKE 'pg_%' AND
                                                               schemaname NOT LIKE 'dbms_%' AND
                                                               schemaname NOT LIKE 'aq$%' AND
                                                               schemaname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                               schemaname NOT LIKE 'utl_%') = 0 THEN 0 ELSE
            SUM(pg_table_size(c.oid)) END AS tables_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname NOT IN ('pg_catalog', 'information_schema', 'sys') AND
                  c.relkind IN ('r', 't','f','p') AND
                  nspname != 'pg_toast')

--this query returns the size of all the non-system indexes in the db
,q  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_indexes WHERE schemaname NOT LIKE 'pg_%' AND
                                                                schemaname NOT LIKE 'dbms_%' AND
                                                                schemaname NOT LIKE 'aq$%' AND
                                                                schemaname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                                schemaname NOT LIKE 'utl_%') = 0 THEN 0 ELSE
            SUM(pg_indexes_size(c.oid)::bigint) END AS indexes_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname NOT IN ('pg_catalog', 'information_schema', 'sys') AND
                  c.relkind IN ('r', 't','f','p') AND
                  nspname != 'pg_toast')

/* this query returns the maximum allowed connection as specified by the dba with the
CREATE DATABASE [...] CONNECTION LIMIT x command. if no limit is specified it returns -1,
else the number of maximum allowed connection is shown*/
,g  AS (SELECT datconnlimit FROM pg_database WHERE datname = (SELECT current_database()))

--this query returns the number of non-system schemas in the db
,i  AS (SELECT COUNT(*) AS schemas_count FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND
                                                                nspname NOT LIKE 'dbms_%' AND
                                                                nspname NOT LIKE 'aq$%' AND
                                                                nspname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                                nspname NOT LIKE 'utl_%')

--this query returns the number of non-system views in the db
,o  AS (SELECT COUNT(*) AS views_count FROM pg_views WHERE schemaname NOT LIKE 'pg_%' AND
                                                          schemaname NOT LIKE 'dbms_%' AND
                                                          schemaname NOT LIKE 'aq$%' AND
                                                          schemaname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                          schemaname NOT LIKE 'utl_%')

--this query returns the name of the db
,z  AS (SELECT current_database() AS db_name)

--this query returns the number of largeobjects in the db
,b  AS (SELECT COUNT(*) AS lobs_count FROM pg_largeobject_metadata)

--this query returns the size of all largeobjects in the db
,y  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_largeobject_metadata) = 0 THEN 0 ELSE
      SUM(OCTET_LENGTH(data)) END AS lobs_size FROM pg_largeobject)

--this query returns the owner of the database
,c  AS (SELECT o.rolname AS db_owner FROM pg_database d JOIN pg_authid o ON (d.datdba = o.oid) WHERE d.datname = current_database())

SELECT * FROM z,c,g,i,a,d,f,e,q,b,y,o;
