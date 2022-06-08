--tested for PostgreSQL and EPAS 10.
--all the data below refer to the db the client is connected to,
--each schema should be analyzed independently

--SET search_path TO <INSERT_SCHEMA_HERE>
SET search_path TO :schema_name;

WITH
--this query returns the name of the schema
z  AS (SELECT CASE
			 WHEN SPLIT_PART(setting,',',1) LIKE '%user%'
				THEN (SELECT CURRENT_USER)
				ELSE SPLIT_PART(setting,',',1)
			END AS schema_name FROM pg_settings WHERE name = 'search_path')

--this query returns the number of all the tables in the schema
,d  AS (SELECT COUNT(*) AS tables_count FROM pg_tables WHERE schemaname = (SELECT schema_name FROM z))

--this query returns the number of all the indexes in the schema
,e  AS (SELECT COUNT(*) AS indexes_count FROM pg_indexes WHERE schemaname = (SELECT schema_name FROM z))

--this query returns the size of all the tables in the schema
,f  AS (SELECT CASE WHEN (SELECT tables_count FROM d) = 0 THEN 0 ELSE
            SUM(pg_table_size(c.oid)) END AS tables_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname = (SELECT schema_name FROM z) AND
                  c.relkind IN ('r', 't','f','p'))

--this query returns the size of all the indexes in the schema
,q  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_indexes WHERE schemaname = (SELECT schema_name FROM z)) = 0 THEN 0 ELSE
            SUM(pg_indexes_size(c.oid)::bigint) END AS indexes_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname = (SELECT schema_name FROM z) AND
                  c.relkind IN ('r', 't','f','p'))
                 
--this query returns the size of all the materialized views in the schema
,s  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_matviews WHERE schemaname = (SELECT schema_name FROM z)) = 0 THEN 0
            ELSE SUM(pg_relation_size(c.oid))
        END AS matviews_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname = (SELECT schema_name FROM z) AND
                  c.relkind = 'm')

--this query returns the size of all the objects in schema
,a AS (SELECT (SELECT tables_size + indexes_size + matviews_size FROM f,q,s) AS schema_size)

--this query returns the number of materialized views in the schema
,n  AS (SELECT COUNT(*) AS matviews_count FROM pg_matviews WHERE schemaname = (SELECT schema_name FROM z))

--this query returns the number of views in the schema
,o  AS (SELECT COUNT(*) AS views_count FROM pg_views WHERE schemaname = (SELECT schema_name FROM z))

--this query returns the owner of the schema
,p  AS (SELECT o.rolname AS schema_owner FROM pg_namespace AS n JOIN pg_authid AS o ON (n.nspowner = o.oid) WHERE n.nspname = (SELECT schema_name FROM z))

SELECT * FROM z,p,d,f,e,q,o,a,n,s

