#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username postgres --dbname postgres <<-EOSQL
    CREATE DATABASE telegram_bot_dev;
    CREATE DATABASE telegram_bot_test;
EOSQL