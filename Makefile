.PHONY: lint lint-fix dev-dep dep test cobertura docker-dep db-migrate db-rollback

lint:
	go fmt ./...
	golangci-lint run --concurrency 2 --color always --timeout 10m0s

lint-fix:
	golangci-lint run --color always --fix

dev-dep:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/dmarkham/enumer@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7

dep:
	go mod tidy
	go mod vendor

docker-dep:
	docker-compose --env-file dev/.env -f dev/docker-compose.yml up --no-recreate

db-migrate:
	atlas migrate apply -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir)

db-rollback:
	atlas migrate down -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir) --to-version $(version) --dev-url "docker://mysql/8/example"
