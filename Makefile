build:
	@go build -o bin/bloq

run: build
	@./bin/bloq

proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

test:
	@go test -v ./...

.PHONY: build run proto test