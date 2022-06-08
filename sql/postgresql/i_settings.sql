--Tested for PostgreSQL 14 and EPAS 14, compatible with PostgreSQL 9 and EPAS 10.
--some settings are unavailable (code -404) for some versions.
WITH
    --returns DB engine
z AS (SELECT CASE
        WHEN (SELECT version()) ILIKE '%enterprisedb%' THEN 'EDB'
        WHEN (SELECT version()) ILIKE '%edb%' THEN 'EDB'
        WHEN (SELECT version()) ILIKE '%postgresql%' THEN 'PostgreSQL'
        ELSE 'Other DB engine' END
        AS db_engine)

    --returns the DB version
,a  AS (WITH aa AS (SELECT REGEXP_MATCHES((SELECT version()),'(([0-9]{1,2}\.){1,3}\S*){1}'))
        SELECT regexp_matches[1] AS version FROM aa) --short version

,za AS (SELECT db_engine||' '||version AS db_version FROM z,a)

    --work_mem setting in a human-readable format
,b  AS (SELECT ((setting::int*1024)::bigint) AS work_mem FROM pg_settings WHERE name = 'work_mem')

    --archive_mode setting
,c  AS (SELECT setting AS archive_mode FROM pg_settings WHERE name = 'archive_mode')

    --archive_command setting
,d  AS (SELECT setting AS archive_command FROM pg_settings WHERE name = 'archive_command')

    --min_wal_size setting in a human-readable format
,e  AS (SELECT CASE
    WHEN (SELECT COUNT(*) FROM pg_settings WHERE name = 'min_wal_size') = 0 THEN '-404'
    ELSE (SELECT setting::int*1024*1024::bigint FROM pg_settings WHERE name = 'min_wal_size') END AS min_wal_size)

    --max_wal_size setting in a human-readable format
,f  AS (SELECT CASE
    WHEN (SELECT COUNT(*) FROM pg_settings WHERE name = 'max_wal_size') = 0 THEN '-404'
    ELSE (SELECT setting::int*1024*1024::bigint FROM pg_settings WHERE name = 'max_wal_size') END AS max_wal_size)

    --max_connections setting
,g  AS (SELECT setting AS max_connections FROM pg_settings WHERE name = 'max_connections')

    --checkpoint_completion_target setting in a human-readable format
,h  AS (SELECT (setting::numeric*100)::smallint||'%' AS checkpoint_completion_target FROM pg_settings WHERE name = 'checkpoint_completion_target')

    --default_statistics_target setting
,i  AS (SELECT setting AS default_statistics_target FROM pg_settings WHERE name = 'default_statistics_target')

    --random_page_cost setting
,j  AS (SELECT setting AS random_page_cost FROM pg_settings WHERE name = 'random_page_cost')

    --maintenance_work_mem setting in a human-readable format
,k  AS (SELECT ((setting::int*1024)::bigint) AS maintenance_work_mem FROM pg_settings WHERE name = 'maintenance_work_mem')

    --shared_buffers setting in a human-readable format
,l  AS (SELECT ((setting::int*1024)::bigint) AS shared_buffers FROM pg_settings WHERE name = 'shared_buffers')

    --effective_cache_size setting in a human-readable format
,m  AS (SELECT ((setting::bigint*(SELECT current_setting('block_size')::bigint))::bigint) AS effective_cache_size FROM pg_settings WHERE name = 'effective_cache_size')

    --effective_io_concurrency setting
,n  AS (SELECT setting AS effective_io_concurrency FROM pg_settings WHERE name = 'effective_io_concurrency')

    --max_worker_processes setting
,o  AS (SELECT CASE
    WHEN (SELECT COUNT(*) FROM pg_settings WHERE name = 'max_worker_processes') = 0 THEN '-404'
    ELSE (SELECT setting FROM pg_settings WHERE name = 'max_worker_processes') END AS max_worker_processes)

    --max_parallel_workers setting
,p  AS (SELECT CASE
    WHEN (SELECT COUNT(*) FROM pg_settings WHERE name = 'max_parallel_workers') = 0 THEN '-404'
    ELSE (SELECT setting FROM pg_settings WHERE name = 'max_parallel_workers') END AS max_parallel_workers)

    --data directory
,r  AS (SELECT setting AS data_dir FROM pg_settings WHERE name = 'data_directory')

    --port
,x  AS (SELECT setting AS port FROM pg_settings WHERE name = 'port')

SELECT * FROM za,r,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p;
