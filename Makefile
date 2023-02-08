tidy:
	go mod tidy && go mod download && go mod vendor

test:
	go test -v -count=1 ./...

lint:
	golangci-lint run -v ./...

lintfix:
	golangci-lint run -v --fix ./...
