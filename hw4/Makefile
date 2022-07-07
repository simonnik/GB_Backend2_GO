up:
	docker-compose --file docker-compose.yaml up -d --remove-orphans

# get service IP-address
# ip address | grep inet

wrk_script:
	wrk \
		-t10 \
		-c200 \
		-d30s \
		-H 'Host: sweet-home.com' \
		--timeout 2s \
		-s ./infrastructure/load_test/wrk_complex.lua \
		http://localhost:8002/entities


wrk_get:
	wrk \
		-t50 \
		-c100 \
		-d10s \
		-H 'Host: sweet-home.com' \
		--timeout 2s \
		http://localhost:8002/entities


# connect to DB in container
# docker exec -ti mysql-container mysql --host=127.0.0.1 --port=3306 -u root -p test
# password = test


# add data to the database (POST)
# curl -X POST -F 'token=admin_secret_token' -F 'id=3' -F 'data=data_three' http://127.0.0.1:8002/entity
# curl -X POST -F 'token=admin_secret_token' -F 'id=4' -F 'data=data_four' http://127.0.0.1:8002/entity


# read entities (GET) in JSON
# curl http://127.0.0.1:8002/entities


# performance test
# wrk -t1 -c1 -d 5m http://localhost:8002/entities
# wrk -t1 -c1 -d5m -s ./load_test/wrk.lua http://localhost:8002