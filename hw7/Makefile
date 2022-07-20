up:
	docker-compose --file docker-compose.yaml up -d --remove-orphans

wrk:
	wrk \
		-t10 \
		-c200 \
		-d30s \
		--timeout 2s \
		-s ./infrastructure/wrk.lua \
		http://localhost:8081
