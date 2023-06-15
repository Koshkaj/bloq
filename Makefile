build:
	go build -o bin/bloq

run: build
	./bin/bloq

test:
	@go test -v ./...