--tested for PostgreSQL and EPAS 10.
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

--this query returns the size of all the materialized views in the db
,s  AS (SELECT CASE WHEN (SELECT COUNT(*) FROM pg_matviews) = 0 THEN 0
            ELSE SUM(pg_total_relation_size(c.oid))
        END AS matviews_size FROM pg_class AS c
            LEFT JOIN pg_namespace AS n ON (n.oid = c.relnamespace)
            WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND
                  c.relkind = 'm' AND
                  nspname != 'pg_toast')

/* this query returns the maximum allowed connection as specified by the dba with the
CREATE DATABASE [...] CONNECTION LIMIT x command. if no limit is specified it returns -1,
else the number of maximum allowed connection is shown*/
,g  AS (SELECT datconnlimit FROM pg_database WHERE datname = (SELECT current_database()))

--this query returns the number of installed extensions
,h  AS (SELECT COUNT(*) AS extensions_count FROM pg_extension)

--this query returns the number of non-system schemas in the db
,i  AS (SELECT COUNT(*) AS schemas_count FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND
                                                                nspname NOT LIKE 'dbms_%' AND
                                                                nspname NOT LIKE 'aq$%' AND
                                                                nspname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                                                nspname NOT LIKE 'utl_%')


/*this query checks if the db has a publication or a subscription (meaning that it is in a logical-replication setup or it once was
and the replication is broken) it returns a boolean value*/
,j  AS (SELECT CASE
                WHEN (SELECT COUNT(*) FROM pg_publication) > 0 OR (SELECT COUNT(*) FROM pg_subscription) > 0 THEN true
                ELSE false END AS logic_repl_setup)

--this query returns the number of publication in the db
,k  AS (SELECT COUNT(*) AS publications_count FROM pg_publication)

--this query returns the number of subscriptions in the db
,l  AS (SELECT COUNT(*) AS subscriptions_count FROM pg_subscription)

--this query returns the number of materialized views in the db
,n  AS (SELECT COUNT(*) AS matviews_count FROM pg_matviews)

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

SELECT * FROM z,g,i,a,d,f,e,q,b,o,n,s,h,j,k,l;