up:
	docker compose up -d --build

down:
	docker compose down -v

seed:
	docker compose exec postgres psql -U postgres -d booking -f /postgres-seed/seed.sql

logs:
	docker compose logs -f app_go

test:
	go test ./... -cover
