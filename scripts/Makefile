up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f $(SERVICE_NAME)

psql:
	docker-compose exec $(SERVICE_NAME) psql -U $(DB_USER) -d $(DB_NAME)

reset:
	docker-compose down -v
	docker-compose up -d
