FROM mcr.microsoft.com/mssql/server:2017-CU24-ubuntu-16.04

WORKDIR /app

# Bundle config source
COPY /entrypoint.sh .
COPY /configure-db.sh .
COPY /schema/setup.sql .
COPY /. .
# Grant permissions for to our scripts to be executable
RUN chmod +x /app/entrypoint.sh
RUN chmod +x /app/configure-db.sh

# entrypoint configuration
ENTRYPOINT ["./entrypoint.sh"]

# cmd to see logs of mssql
CMD ["tail -f /dev/null"]