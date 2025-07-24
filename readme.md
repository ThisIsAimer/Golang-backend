# my simple backend project with golang
## common go mod cmds:-

- go mod init |module name|
- go get |package link|
- go mod tidy (removes all the unnecessary packages not in use in the go mod file)
- go list -m all
- go clean -modcache
- go mod verify

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


# middleware snippits

            "func MiddlewareName(next http.Handler) http.Handler{",
			"",
			"    return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request)  {",
			"        next.ServeHTTP(w,r)",
			"    }),",
			"",
			"}",
