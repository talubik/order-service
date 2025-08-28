export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=talub
export DB_PASSWORD=boba
export DB_NAME=mydatabase
run:
	docker compose -f docker-compose.postgres.yml up -d && \
	docker compose -f docker-compose.kafka.yml up -d &&\
	go run src/cmd/service/main.go
test:
	go test ./...
clean:
	docker compose -f docker-compose.postgres.yml down -v && \
	docker compose -f docker-compose.kafka.yml down -v 