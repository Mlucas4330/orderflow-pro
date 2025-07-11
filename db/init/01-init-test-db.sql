SELECT 'CREATE DATABASE orderflow_pro_test'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'orderflow_pro_test')\gexec