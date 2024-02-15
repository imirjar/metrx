## Создание базы данных 

Вот пример создания БД PostgreSQL на Linux:

```
sudo -i -u postgres
```
```
psql -U postgres
```
```
create database metrics;
```
```
create user metrics_master with encrypted password 'userpassword';
```
```
grant all privileges on database dbname to metrics_master; 
```

Для создания таблицы metrics выполните команду

```
CREATE TABLE metrics (`"id" VARCHAR(50) NOT NULL, "type" VARCHAR(250) NOT NULL, "delta" FLOAT, "value" INTEGER`) 
```