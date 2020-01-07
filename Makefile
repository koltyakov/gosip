install:
	go get -u ./... && go mod tidy

test:
	go test ./... -v -race -count=1

test-api:
	SPAUTH_ENVCODE=spo go test ./api/... -v -count=1

format:
	gofmt -s -w .

coverage:
	SPAUTH_ENVCODE=spo SPAPI_HEAVY_TESTS=true go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out