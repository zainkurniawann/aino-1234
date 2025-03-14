# Backend

This project is a **API** built with **Golang** and **Echo framework**.

## Installation

### Requirements  
- [Golang 1.23.3](https://go.dev/doc/install)  
- [Postgresql 16.3](https://www.postgresql.org/download/)  

### Building  
```sh
git clone https://github.com/zainkurniawann/aino-1234.git
cd aino_doc
go run .
```

## Usage
Base URL = http://:1234

### Run this first
```sh
cd aino_doc
go run .\database\migrations\goose.go
```
or manually with
```sh
// up for migrate (all)
goose -dir ./databases/migrations postgres "user=YOUR_POSTGRES_USER password=YOUR_POSTGRES_PASSWORD dbname=YOUR_POSTGRES_DB sslmode=disable" up
// down for rollback (once)
goose -dir ./databases/migrations postgres "user=YOUR_POSTGRES_USER password=YOUR_POSTGRES_PASSWORD dbname=YOUR_POSTGRES_DB sslmode=disable" down
```
#### Notes
if you use this
```sh
cd aino_doc
go run .\database\migrations\goose.go
```
make sure you have configured the database/connection.go file

## list all form
// You can change it according to your data stored in document_ms.
- Dampak Analisa -> filtered by document_code = 'DA' // You can change it in service.go.
- ITCM -> filtered by document_code = 'ITCM'
- Berita Acara -> filtered by document_code = 'BA'
- Hak Akses -> filtered by document_code = 'HA'

## auth
// You can change it in middleware, according to your data stored in role_ms.
Role required:
- Member -> middleware with role_code = 'M' 
- Admin -> middleware with role_code = 'A'
- Superadmin -> middleware with role_code = 'SA'

### Login
In database migration there is already a user that can be used to log in. please check users table