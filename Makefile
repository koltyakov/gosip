install:
	go get -u ./... && go mod tidy

test:
	go test ./... -v -race -count=1

format:
	gofmt -s -w .
