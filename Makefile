.PHONY: lint lint-fix dev-dep dep test cobertura docker-dep db-migrate db-rollback

docker-dep:
	docker-compose --env-file dev/.env -f dev/docker-compose.yml up --no-recreate

db-migrate:
	atlas migrate apply -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir)

db-rollback:
	atlas migrate down -u "mysql://$(db_username):$(db_password)@$(db_host):$(db_port)/$(db_name)" --dir file://$(migration_dir) --to-version $(version) --dev-url "docker://mysql/8/example"
