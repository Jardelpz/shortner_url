up:
	docker-compose up -d

down:
	docker-compose down

reset:
	docker-compose down -v
	docker-compose up -d

test:
	go test ./...

test/integration:
	go test -tags=integration ./...
