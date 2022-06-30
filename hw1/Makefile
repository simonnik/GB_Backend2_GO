# start postgres container
.PHONY: START_PG

START_PG:
	docker-compose up -d

# psql connection string (limited user):
# 	export PGPASSWORD='appp@$$w0rd'; psql -h 127.0.0.1 -p 5432 -U appuser -d app_db

# connect to container
# docker exec -ti pg_container psql -U appuser app_db