# payment-practice

API with functionality of Payment and Managing accounts.

### Prerequisites:
```
Go 1.22
Postgresql
Docker-compose
```

### How to use:

```
Steps:
- make build
- make run

This command will use your database for build project with this local variables and build docker image.

- Turn off you database container
- Run docker-compose up -d
- migrate the database in migration file [doc](https://github.com/arthTes/payment-api/blob/main/scripts/database/migrations/init_db.sql)

Use the postman collection for test

service running in :8080

```

### Documentation

API Documentation in :
[doc](https://github.com/arthTes/payment-api/tree/main/docs/payment-api.yaml)

### Future Work
```
-Dockerfile compound
``` 