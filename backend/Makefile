all:
	go build -o bin/createuser ./cmd/createuser/
	go build -o bin/localserver ./cmd/localserver/
	chmod +x ./bin/*
	@echo "Build done"

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html