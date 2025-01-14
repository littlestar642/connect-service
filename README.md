# Counter-service

## Description
The service used to count number of requests, every minute. 

## Prerequisites

- Golang - ^1.23.2
- Docker - ^27.2.0
- Docker Compose - ^2.29.2

If not using docker

- Redis-cli
- kafka

## Build Instructions 

```shell script
git clone git@github.com:littlestar642/counter-service.git

cd counter-service
```

Local
```shell script
mkdir build

go build -o build/counter-service  cmd/counter.go
```

Docker
```shell script
make build
```

## Run Insructions

Local

- might need to udpate kafka and redis address

```shell script
./build/counter-service 
```

Docker
```shell script
make compose-run
```

## APIs

GET - Accept request id (optionally endpoint) and increment count

```shell script
curl http://localhost:8080/api/verve/accept?id=0
```

POST - Accept request count
```shell script
curl -X POST http://localhost:8080/api/verve/accept?count=0
```




