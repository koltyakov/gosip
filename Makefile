install:
	go get -u ./... && go mod tidy

test:
	go test ./... -v -race -count=1

test-api:
	SPAUTH_ENVCODE=spo go test ./api/... -v -count=1

format:
	gofmt -s -w .
