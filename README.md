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