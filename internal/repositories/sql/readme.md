this is now you make an sql query

```sql
create database if not exists  (DATABASE NAME);
use (DATABASE NAME);
create table if not exists (TABLE NAME)(
 id int auto_increment primary key,
 first_name varchar(255) not null,
 last_name varchar(255),
 email varchar(255) not null unique,
 class varchar(255) not null,
 subject varchar(255) not null,
 index(email)
)
```