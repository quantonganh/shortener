frontend:
	cd ui && npm install && npm run build

lint:
	golangci-lint run -v ./...

test:
	go test -v ./...

integration-test:
	go test -tags integration -v -run TestShortener ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -ldflags="-s -w" -o shortener cmd/shortener/main.go