all:
	go build -o bin/charon .
	@echo "Build done"

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
