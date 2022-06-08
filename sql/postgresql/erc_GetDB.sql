SELECT datname FROM pg_database WHERE datname NOT IN ('template1', 'template0');
