# loan-backend

## Service Architecture Overview

Database Design
https://dbdiagram.io/d/Amartha-Loan-69c95965fb2db18e3b2e3c69

## Development Guide

### Prerequisite

- [Atlas](https://atlasgo.io/docs#installation)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Setup
- Spin up dependencies with docker-compose
  ```
  make docker-dep
  ```

- Migrate the database
  ```sh
  make db-migrate
  ```

### Local development using docker-compose

We provide `docker-compose.yaml` file to hold the non service dependencies of backend such as MySQL. To spin up the dependencies, you can do the following steps.

- Copy `dev/env.sample` to `dev/.env` and if necessary, modify the env value(s)
- Spin up the dependencies.

  ```sh
  make docker-dep
  ```

### Database Migration

- Create new migration
  ```sh
  update schema.sql
  atlas migrate diff {migration_filename} --to file://{schema_path} --dev-url "docker://mysql/8/dev"
  ```

- Up migration
  ```sh
  atlas migrate apply -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir)
  ```

- Down migration

  ```sh
  atlas migrate down -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir) --to-version $(version) --dev-url "docker://mysql/8/example"
  ```

### Generate Mock

example generate mock Interface
```
mockgen -destination=mock/user_repository.go -package=mock example.com/loan/module/user/internal/repository UserRepository

mockgen -destination=mock/user_usecase.go -package=mock example.com/loan/module/user/internal/usecase UserUsecase
```

## Program Simulation

### Setup

Make sure that you already did spin up the dependencies with docker-compose and did the database migration.

### Create User, Account, Loan and Force Disburse Loan

Create User with Username testuser2 and Status active
```
go run cmd/cli/main.go create-user --username=testuser2 --status=active
```

Create Account with User ID 2, Balance 1 Mio IDR and Status active
```
go run cmd/cli/main.go create-account --user_id=2 --balance=1000000 --status=active
```

Create Loan with Principal 5 Mio IDR with Term 50 weeks, with Interest 10%
```
go run cmd/cli/main.go create-loan --user_id=2 --principal=5000000 --term=50 --interest=0.1
```

Force Disburse Loan with Loan ID 4, to trigger bulk create Loan Billing
```
go run cmd/cli/main.go force-disburse-loan --loan_id=4
```

### Get Outstanding for Certain Loan through endpoint

Hit localhost:8888/loans/{id}/outstanding

### Get User Status to Check IsDelinquent through endpoint

Hit localhost:8888/users/{id} and check status column

### Make Payment through Endpoint

Hit POST localhost:8888/loans/billings/{id}/pay and make payment for that loan billing
