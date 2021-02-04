BINARY=svc

dev:
	@go mod tidy
	@go mod vendor
	@env GO111MODULE=on CGO_ENABLED=0 go build -o ${BINARY} -mod=vendor

build:
	env GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY} -mod=vendor -a -installsuffix cgo -ldflags '-w'

docker:
	@echo "> Build Docker image sicepat/svc"
	@docker build -t sicepat/dsvc -f Dockerfile .
	@docker build -t sicepat/svc -f Dockerfile.dev . 

run: 
	@./svc 

run-dev: dev run


.PHONY: build deps dev docker
