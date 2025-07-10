# my simple backend project with golang
## common go mod cmds:-

- go mod init |module name|
- go get |package link|
- go mod tidy (removes all the unnecessary packages not in use in the go mod file)
- go list -m all

#### module is a collection of related packages

## https cert cmd
#### download from https://slproweb.com/products/Win32OpenSSL.html
#### then put bin folder in env variables

- openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out certificate.pem -days 365