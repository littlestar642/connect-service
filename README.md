# Counter-service

## Description
The service used to count number of requests, every minute. 

## Prerequisites

- Go - ^1.23.2
- Docker - ^27.2.0
- Docker Compose - ^2.29.2
- Make

If not using docker

- Redis-cli
- kafka

## Build Instructions 

```shell script
git clone git@github.com:littlestar642/counter-service.git

cd counter-service
```

Docker
```shell script
make docker-build
```

Local
```shell script
mkdir build

go build -o build/counter-service  cmd/counter.go
```

## Run Insructions

Docker
```shell script
make compose-run
```

Local

- might need to udpate kafka and redis address

```shell script
./build/counter-service 
```

## APIs

GET - Accept request id (and optionally endpoint) to increment count

```shell script
curl http://localhost:8080/api/verve/accept?id=0
```

```shell script
curl http://localhost:8080/api/verve/accept?id=43&endpoint=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fverve%2Faccept
```

Response 

- "ok"     - denotes success
- "failed" - denotes failure


For testing purpose:

POST - Accept request count
```shell script
curl -X POST http://localhost:8080/api/verve/accept?count=0
```

[Thought Process](./thought-process.md)
