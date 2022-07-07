ARGS=$(filter-out $@, $(MAKECMDGOALS))

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)

PROJECT = github.com/simonnik/GB_Backend1_CW_GO
GO_PKG = $(shell go list $(PROJECT)/...)

POSTGRESQL_URL=postgresql://$$APP_DB_USER:$$APP_DB_PASSWORD@$$APP_DB_HOST:$$APP_DB_PORT/$$APP_DB_NAME?sslmode=$$APP_DB_SSL_MODE
MGR_COMMAND=-path=$$APP_DB_MIGRATIONS_DIR -database $(POSTGRESQL_URL) -verbose

bootstrap:
ifndef HAS_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
endif
ifndef HAS_IMPORTS
	go mod download golang.org/x/tools
	go mod tidy
endif

up: docker-up
down: docker-down
stop: docker-stop
start: docker-start
restart: docker-restart
build: docker-build
init: bootstrap \
	  pre-commit-install \
	  create-env \
	  docker-down-clear \
	  docker-pull \
	  docker-build \
	  docker-up \
	  migrate-up

create-env:
	if [ ! -f './.env' ]; then cp ./.env.sample ./.env; else exit 0; fi;
	#cp -n .env.sample .env && exit 0

docker-up:
	docker-compose up --detach --remove-orphans

docker-down:
	docker-compose down --remove-orphans

docker-down-clear:
	docker-compose down --volumes --remove-orphans

docker-stop:
	docker-compose stop

docker-start:
	docker-compose start

docker-restart:
	docker-compose restart

docker-pull:
	docker-compose pull

docker-build:
	docker-compose build

docker-img-build-shortener:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -f build/app/Dockerfile -t	sim4nix/shortener:$(ARGS) .
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -f build/migrate/Dockerfile -t	sim4nix/shortener-migrate:$(ARGS) .
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -f build/swagger/Dockerfile -t	sim4nix/shortener-swagger:$(ARGS) .

docker-img-build-migrate:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -f build/migrate/Dockerfile -t	sim4nix/shortener-migrate:$(ARGS) .

docker-img-build-swagger:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -f build/swagger/Dockerfile -t	sim4nix/shortener-swagger:$(ARGS) .

docker-img-push-shortener:
	docker push sim4nix/shortener:$(ARGS)

docker-img-push-migrate:
	docker push sim4nix/shortener-migrate:$(ARGS)

docker-img-push-swagger:
	docker push sim4nix/shortener-swagger:$(ARGS)

kub-roll:
	kubectl -n g3-snew apply -f manifests-roll.yml

kub-fixed:
	kubectl -n g3-snew apply -f manifests-fixed.yml

go-build: pre-commit test
	go build -o bin/main ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

migrate-up:
	export `grep -v "^#" .env | xargs` && \
	env docker-compose run --rm migrate $(MGR_COMMAND) up

migrate-create:
	export `grep -v "^#" .env | xargs` && \
	env docker-compose run --rm migrate $(MGR_COMMAND) create -ext sql -dir $$APP_DB_MIGRATIONS_DIR -seq $(ARGS)
	sudo chown -R $$USER:$$GROUP ./migrations

migrate-down:
	export `grep -v "^#" .env | xargs` && \
	env docker-compose run --rm migrate $(MGR_COMMAND) down 1

migrate-force:
	export `grep -v "^#" .env | xargs` && \
	env docker-compose run --rm migrate $(MGR_COMMAND) force $(ARGS)

test:
	@echo "+ $@"
	@go list -f '"go test -cover -v {{.Dir}}"' $(GO_PKG) | xargs -L 1 sh -c

fmt:
	@echo "+ $@"
	@go list -f '"gofmt -w -s -l {{.Dir}}"' $(GO_PKG) | xargs -L 1 sh -c

imports:
	@echo "+ $@"
	@go list -f '"goimports -w {{.Dir}}"' ${GO_PKG} | xargs -L 1 sh -c

check: bootstrap
	@echo "+ $@"
	@golangci-lint run ./...

pre-commit:
	@echo "+ $@"
	pre-commit run --all-files

pre-commit-install:
	@echo "+ $@"
	pre-commit install

%:
	@true

.PHONY: bootstrap \
	build \
	test \
	fmt \
	imports \
	check \
	init \
	run