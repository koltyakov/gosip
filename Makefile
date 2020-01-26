install:
	go get -u ./... && go mod tidy

test-auth:
	go test ./... -v -race -count=1

test-api:
	SPAUTH_ENVCODE=spo go test ./api -v -count=1
	SPAUTH_ENVCODE=spo go test ./api/csom -v -count=1

format:
	gofmt -s -w .

generate:
	go install ./cmd/ggen/...
	go generate ./api/...
	make format

coverage:
	SPAUTH_ENVCODE=spo SPAPI_HEAVY_TESTS=true go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

auth-precover:
	bash ./test/scripts/cover-auth.sh