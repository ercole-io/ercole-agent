SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%'   AND
                                       nspname NOT LIKE 'dbms_%' AND
                                       nspname NOT LIKE 'aq$%'   AND
                                       nspname NOT IN ('information_schema', 'sys','msg_prop_t', 'dbo') AND
                                       nspname NOT LIKE 'utl_%';