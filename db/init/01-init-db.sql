GRANT ALL PRIVILEGES ON DATABASE postgres TO :"user";

SELECT 'CREATE DATABASE orderflow_dev_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'orderflow_dev_db')\gexec

SELECT 'CREATE DATABASE orderflow_test_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'orderflow_test_db')\gexec