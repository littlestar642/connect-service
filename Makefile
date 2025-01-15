docker-build:
	go fmt ./...
	docker build -t counter-service-app:latest .

compose-run:
	$(MAKE) docker-build
	docker-compose up

build-binary:
	go build -o build/counter-service cmd/counter.go

run:
	go run cmd/counter.go

test:
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out