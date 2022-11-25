Gin-N-Juice Framework
=================

Gin-N-Juice is a Go JSON API framework that is built on [Gin](https://github.com/gin-gonic/gin),
[Gorm](https://gorm.io) and [Goose](https://github.com/pressly/goose). 

Current Features
------
- Uses Gorm as an ORM
- Database migrations via Goose
- Simple routing via Gin
- Simple Auth (login,register,email verification,forgot password and reset password)
- File watcher to auto reboot server during development

Requirements
------
- The only mail provider that currently works is Mailgun.
- Go version 1.18 or later

Setup
------
Create a `.env` in the root directory of the project. The required environment variables are:
```
ENCRYPT_KEY=[randomstring]
MAILGUN_DOMAIN=[get from mailgun]
MAILGUN_PRIVATE_KEY=[get from mailgun]
MAILGUN_VALIDATION_KEY=[get from mailgun]
DB_TYPE=[sqlite,postgres or mysql]
DB_CONNECTION_STRING="[connection string for your db, see below]"
PORT=[port to run server on. defaults to port 8080]
SERVER_HOSTNAME=[hostname server will listen for. defaults to 127.0.0.1/localhost]
MODE=[set to "production", otherwise runs in debug mode]
EMAIL_FROM=[default email address to send messages from]
PACKAGE_NAME=[if you change your package name in go.mod update this to match]
```
Example DB connection string reference: 
```
MySQL: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
SQLite: "gorm.db"
PostgreSQL: "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=America/New_York"
```
Add routes to the `routes` directory and reference them in the `routes/router.go`.

Commands
------
- `go run .` or `go run . serve`
  - Runs the webserver
  - Also sets up a watcher and will auto restart webserver when files are changed
- `go run . migrate [status|up|down|create|etc...]`
  - Runs goose migrations, see goose's documentation for more details
- `go run . rename [package-name]`
  - If you want to change your package name from `gin-n-juice` to something else. It will update the
    `go.mod` and all `.go` file imports
- `go run . test`
  - Run tests on all routes and models
- `go run . generator -type [model|resource] -name [ex:user] -fields [name:type,name:type,...]`
  - `-fields` is only used with model type
  - resource creates standard API endpoints for the name (create,update,delete,list,view)
  - use singular names, the name will automatically be pluralized when needed (example: `user` not `users`)
  - Example: `go run . generator -type model -name book -fields title:string,category_id:uint,pages:uint`
    - Creates a single file `book.go` in the `/models` directory
  - Example: `go run . generator -type resource -name book`
    - Creates a new directory `/routes/books`
    - Creates five files in that directory
      - `create.go`, `delete.go`, `get.go`, `list.go`, `update.go`
