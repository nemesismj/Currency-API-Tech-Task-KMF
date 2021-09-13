## README

#### Table of contents

- [Docker](#Docker)
- [Server](#Server)
- [MSSQL](#MSSQL)
- [TESTS](#TESTS)
- [Swagger](#Swagger)


### Docker

#### Description

Чтобы поднять базу данных в контейнере докер нужно

```
1) docker-compose build
2) docker-compose up
```

### Server
#### DESCRIPTION
Для запуска сервера нужно сначала нужно выполнить шаг с докером после находясь в директории проекта выполнить команду
```
1) make
2) ./apiserver
```
После у вас должно появится что сервер был запущен

#### DESCRIPTION

```
{
"host": "localhost",
"port": "8080"
}
```


### MSSQL

#### DESCRIPTION
```
{
"user": "kursUser",
"pswd": "kursPswd",
"connection_string": "sqlserver://kursUser:kursPswd@localhost:1400/?charset=utf8"
}
```

### TESTS

#### Description

Чтобы запустить тест, нужно запустить команду make test


### Swagger

#### Description
```
документация swagger доступна по пути -> http://localhost:8080/swagger/index.html/
```

