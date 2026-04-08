# loan-backend

## Service Architecture Overview

Database Design
https://dbdiagram.io/d/Amartha-Loan-69c95965fb2db18e3b2e3c69

## Development Guide

### Prerequisite

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Go 1.24.1 or later](https://golang.org/doc/install)
- [Atlas](https://atlasgo.io/docs#installation)
- [Docker](https://docs.docker.com/engine/install/)
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

Overdue Billing Checker to trigger update billing to overdue if the condition match and user status to delinquent if condition match
```
go run cmd/cli/main.go overdue-billing-checker
```

### Get Outstanding for Certain Loan through endpoint

After creating loan and force disburse it, you could see the outstanding bill for that loan

Hit localhost:8888/loans/{id}/outstanding

<img width="2944" height="1838" alt="image" src="https://github.com/user-attachments/assets/7070c1fc-6da3-4619-84e7-b85d5b9fa7f9" />

### Make Payment through Endpoint

After Creating loan and force disburse it, you could pay the billing

Hit POST localhost:8888/loans/billings/{id}/pay and make payment for that loan billing

<img width="2944" height="1840" alt="image" src="https://github.com/user-attachments/assets/a9e9b249-407d-424d-8a69-df51119273ac" />

### Get User Status to Check IsDelinquent through endpoint

To see the delinquent status, you need to create the billings through force disburse.
After that, change the due date for at least 2 billings before today.
Run script/command overdue-billing-checker to trigger user status change to delinquent.

Hit localhost:8888/users/{id} and check status column (1: inactive, 2: active, 3: delinquent)

<img width="2942" height="1824" alt="image" src="https://github.com/user-attachments/assets/a930dbf0-862b-4731-8cf1-ee5dd162d332" />
