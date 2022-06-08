--tested for PostgreSQL 9 and EPAS 10
WITH
    --this column shows the maximum allowed connections to the DB
 a  AS (SELECT setting AS max_connections FROM pg_settings WHERE name = 'max_connections')

    /*this column shows the size of the instance as a sum of the size of all the databases
    as reported in the pg_database system table. a human-readable format is used*/
,q  AS (WITH totals AS (select datname, pg_database_size(datname) AS db_size FROM pg_database)
                SELECT SUM(db_size) AS instance_size FROM totals)

    --this column shows the number of users in the instance. it also counts the default user[s]
,w  AS (SELECT COUNT(*) AS users_num FROM pg_user)

    --this column shows the number of databases in the instance. it does not count the default databases
,x  AS (SELECT COUNT(*) AS db_num FROM pg_database WHERE datname != 'postgres' AND datname NOT LIKE 'template%')

    /*this column shows the number of tablespaces in the instance.
    it does not count the default '$PGDATA' directory as a tablespace*/
,y  AS (SELECT COUNT(*) AS tblsp_num FROM pg_tablespace WHERE spcname NOT LIKE 'pg_%')

SELECT * FROM a,q,w,x,y;
