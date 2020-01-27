install:
	go get -u ./... && go mod tidy

test-auth:
	mkdir -p tmp
	go test ./auth/... -v -race -count=1 -coverprofile=./tmp/auth_coverage.out

test-utils:
	mkdir -p tmp
	SPAUTH_ENVCODE=spo go test ./csom -v -race -count=1 -coverprofile=./tmp/csom_coverage.out
	SPAUTH_ENVCODE=spo go test ./cpass -v -race -count=1 -coverprofile=./tmp/cpass_coverage.out

test-api:
	mkdir -p tmp
	SPAUTH_ENVCODE=spo go test ./api/... -v -count=1 -coverprofile=./tmp/api_coverage.out

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