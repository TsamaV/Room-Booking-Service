up:
	docker-compose up -d --build

down:
	docker-compose down

build:
	go build -o bin/app ./cmd/api/

run:
	go run ./cmd/api

migrate:
	go run ./migrations/auto.go

test:
	go test ./...

test-unit:
	go test ./internal/... ./pkg/... ./configs/...

test-integration:
	go test ./... -run TestCreateRoomScheduleBooking -v
	go test ./... -run TestCancelBooking -v