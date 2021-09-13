CREATE DATABASE $(MSSQL_DB);
GO
USE $(MSSQL_DB);
GO
CREATE LOGIN $(MSSQL_USER) WITH PASSWORD = '$(MSSQL_PASSWORD)', CHECK_POLICY = OFF;
GO
CREATE USER $(MSSQL_USER) FOR LOGIN $(MSSQL_USER);
GO
ALTER SERVER ROLE sysadmin ADD MEMBER [$(MSSQL_USER)];
GO
USE $(MSSQL_DB);
GO
CREATE TABLE R_CURRENCY
(
    id     int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    title  nvarchar(60) not null,
    code   varchar(3) not null,
    value  numeric(18, 2) not null,
    a_date date not null
)
GO