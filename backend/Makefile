targetos=windows linux darwin
targetarch=386 amd64

all:
	go build -o bin/createuser ./cmd/createuser/
	go build -o bin/localserver ./cmd/localserver/
	go build -o bin/centralserver ./cmd/centralserver/
	chmod +x ./bin/*
	@echo "Build done"

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

decrypt:
	rm ./bin/timelock_*
	rm ./bin/timelock.zip
	$(foreach os, $(targetos), \
		$(foreach arch, $(targetarch), \
			GOOS=$(os) GOARCH=$(arch) go build -o bin/timelock_$(os)_$(arch) ./cmd/decrypt; \
		) \
	)
	$(foreach arch, $(targetarch), \
		mv ./bin/timelock_windows_$(arch) ./bin/timelock_windows_$(arch).exe; \
	)
	zip ./bin/timelock.zip ./bin/timelock_*
