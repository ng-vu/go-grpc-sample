# Go starter project

This repository contains sample code for a service called **Blue**.

There are:
- AgencyStaff
- Partner (ServiceProvider)
- Admin

## Quick Start

- Install Go
- Install [glide](https://github.com/Masterminds/glide)
- Install [protobuf](https://github.com/google/protobuf/releases/tag/v3.3.0)
- Put the project at `$GOPATH/src/github.com/ng-vu/go-grpc-sample`

### 1. Install dependencies and tools

```bash
$ cd $GOPATH/src/github.com/ng-vu/go-grpc-sample
$ glide install
$ go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/{protoc-gen-grpc-gateway,protoc-gen-swagger}
$ go install ./vendor/github.com/golang/protobuf/{proto,protoc-gen-go}
```

### 2. Setup Docker

```bash
$ cd $GOPATH/src/github.com/ng-vu/go-grpc-sample
$ docker volume create --name=blue_mongo_data
$ docker volume create --name=blue_postgres_data
$ docker-compose up -d
```

Edit `/etc/hosts` to include this line

```
127.0.0.1    redis postgres mongo
```

### 3. Build and start

```
$ go install github.com/ng-vu/go-grpc-sample/blue/cmd/blue_service
$ $GOPATH/bin/blue_service
```

Those endpoints will be avalable in HTTP or GRPC. Open them in browser to view API documentation.

| Address | Protocol | Service |
|-|-|-|
| [localhost:3000](http://localhost:3000) | HTTP | Agency
|  localhost:3001                         | GRPC | Agency
| [localhost:3100](http://localhost:3100) | HTTP | Partner
|  localhost:3101                         | GRPC | Partner
| [localhost:3900](http://localhost:3900) | HTTP | SAdmin
|  localhost:3901                         | GRPC | SAdmin

### 4. Run tests

```bash
$ cd $GOPATH/src/github.com/ng-vu/go-grpc-sample
$ docker-compose up -d

$ go test -v github.com/ng-vu/go-grpc-sample/blue/...
```

## Try with sample data

### 1. Create sample service provider and agency staff

```
$ go run ./scripts/seed_data/seed_data.go
```

The above command with output `APIKEY`. Remember it to create sample orders.

### 2. Create sample orders

```bash
# Replace <APIKEY> with `APIKEY` returned from running seed_data.go
$ go run ./scripts/seed_orders/seed_orders.go --apikey <APIKEY>
```

### 3. Login and retrieve orders

```bash
$ curl -X POST http://localhost:3000/api/v1/account/login \
  -H 'content-type: application/json' \
  -d '{ "phone": "0123456789", "password": "123456" }'
```

The above command will output `access_token`. Remember it for subsequent requests.

```bash
# Replace <ACCESS_TOKEN> with `access_token` returned from account/login
$ curl -X POST http://localhost:3000/api/v1/customer/lookup \
  -H 'authorization: Bearer <ACCESS_TOKEN>' \
  -H 'content-type: application/json' \
  -d '{	"phone": "01234560002" }'
```

## Code generation

```bash
# Generate protobuf and swagger documents
$ ./scripts/protobuf-gen.sh

# Generate ID constants
$ go run ./scripts/idconsts_gen/idconsts_gen.go
```
