
build:
	go build -o application cmd/main.go
run:
	go run cmd/main.go
test:
	go test ./...