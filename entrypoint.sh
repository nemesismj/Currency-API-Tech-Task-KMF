#!/bin/bash

# Start SQL Server
/opt/mssql/bin/sqlservr &

# Start the script to create the DB and user
/app/configure-db.sh

#/usr/config/runapp.sh

# Call extra command
eval $1