build:
	go fmt ./...
	docker build -t counter-service-app:latest .

compose-run:
	$(MAKE) build
	docker-compose up

run:
	go run cmd/counter.go

test:
	go test ./... -v --cover

jstypes:
	go run ./plugins/jsvm/internal/types/types.go

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out