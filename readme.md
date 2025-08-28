# my simple backend project with golang
## common go mod cmds:-

- go mod init |module name|
- go get |package link|
- go mod tidy (removes all the unnecessary packages not in use in the go mod file)
- go list -m all
- go clean -modcache
- go mod verify

### to use wsl (wsl -d ubuntu)

### imports used!
- go driver for my sql (github.com/go-sql-driver/mysql)
- for reading env file (github.com/joho/godotenv)

#### module is a collection of related packages

## https cert cmd
#### download from https://slproweb.com/products/Win32OpenSSL.html
### then put bin folder in env variables

#### without config file
- openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out certificate.pem -days 365
#### with config file
- openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out certificate.pem -days 365 -config openssl.cnf

#### for testing add  the self signed certificate to postman as authorised certs


## Api testing

- postman
- curl (curl -v -k https://localhost:3000/)


## benchmarking
### use of wsl

- wrk -t8 -c400 -d30s "link"
- h2load -n  1000 -c 100 -t 8 (--h1 for http1)

## For email testing (mailhog)

- download link: https://github.com/mailhog/MailHog
- port http://localhost:8025/


# middleware snippits
```go
			"func MiddlewareName(next http.Handler) http.Handler{",
			"",
			"    return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request)  {",
			"        next.ServeHTTP(w,r)",
			"    })",
			"",
			"}",
```

# CRUD operations in SQL

```sql
- CREATE DATABASE test_database; | (create a new database)
- SHOW DATABASES; | (shows all databases)
- USE test_database; | (tells mysql to use test_database)
- CREATE TABLE my_table( id INT PRIMARY_KEY AUTO_INCREMENT, name VARCHAR(50), age INT); | (set a table in database)
- INSERT INTO my_table(name, age) VALUES("Tanjiro Kamado", 16),("Nezuko Kamado", 14); | (put values in the table)
- SELECT * FROM my_table WHERE id = 2; | (see values from the table)
- UPDATE my_table SET age = 19 WHERE id = 1; | (update existing values in the table)
- DELETE FROM my_table WHERE id = 4; | (delete a row in the table)
- RENAME TABLE my_table to new_table; | (renames table)
- DROP TABLE new_table; / DROP DATABASE Test_database; | (deletes table or database)

- -- for relational databasing
- DROP INDEX IF EXISTS idx_class ON teachers; | (for droping index)
- CREATE INDEX idx_class ON teachers(class); | (creates a foreign index key for relational database)
- FOREIGN KEY (class) REFERENCES teachers(class) | (links foreign key to column)

- --For time operations
- user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```