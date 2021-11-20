test:
# write tests and uncomment
#   go test ./...

setup: test
	$(eval BINARY=gnore)

build: setup
	go build -o ${BINARY} ./cmd/gnore/main.go

run:
	go run ./cmd/gnore/main.go