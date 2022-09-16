Gin-N-Juice Framework
=================

Gin-N-Juice is a Go JSON API framework that is built on [Gin](https://github.com/gin-gonic/gin),
[Gorm](https://gorm.io) and [Goose](https://github.com/pressly/goose). 

Current Features
------
- Uses Gorm as an ORM
- Database migrations via Goose
- Simple routing via Gin
- Simple Auth (login,register,forgot password and reset password)

Requirements
------
The only mail provider that currently works is Mailgun.

Setup
------
Create a `.env` in the root directory of the project. The required environment variables are:
```
ENCRYPT_KEY=[randomstring]
MAILGUN_DOMAIN=[get from mailgun]
MAILGUN_PRIVATE_KEY=[get from mailgun]
MAILGUN_VALIDATION_KEY=[get from mailgun]
DB_TYPE=[sqlite,postgres or mysql]
DB_CONNECTION_STRING="[connection string for your db]"
```
Add routes to the `routes` directory and reference them in the `routes/router.go`.

Commands
------
- `go run .` or `go run . serve`
  - Runs the webserver
- `go run . migrate [status|up|down|etc...]`
  - Runs goose migrations, see goose's documentation for more details
- `go run . test`
  - Runs tests on all routes currently, plan on testing models also
