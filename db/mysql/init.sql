CREATE SCHEMA IF NOT EXISTS redwing;

CREATE USER IF NOT EXISTS 'redwing'@'%' IDENTIFIED WITH 'caching_sha2_password'
        BY 'redwing';

GRANT ALL ON redwing.* TO 'redwing'@'%';