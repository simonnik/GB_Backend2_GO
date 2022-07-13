up:
	docker-compose --file docker-compose.yaml up -d --remove-orphans

migrate-up:
	docker-compose run --rm -v "${PWD}/migrations:/app/migrations" goose /app/goose -dir /app/migrations postgres 'host=shard_0	user=test password=test dbname=test sslmode=disable' up
	docker-compose run --rm -v "${PWD}/migrations:/app/migrations" goose -- -dir /app/migrations postgres 'host=shard_1 user=test password=test dbname=test sslmode=disable' up
	docker-compose run --rm -v "${PWD}/migrations:/app/migrations" goose -- -dir /app/migrations postgres 'host=shard_2 user=test password=test dbname=test sslmode=disable' up

fill-users:
	go run ./cmd/fill/users

fill-activities:
	go run ./cmd/fill/activities