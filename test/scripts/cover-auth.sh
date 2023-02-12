SPAPI_SKIP_TESTS=true

go test ./auth/azurecert/... -coverprofile=auth/azurecert/coverage.data -covermode=atomic
go test ./auth/azurecreds/... -coverprofile=auth/azurecreds/coverage.data -covermode=atomic
go test ./auth/azureenv/... -coverprofile=auth/azureenv/coverage.data -covermode=atomic
go test ./auth/device/... -coverprofile=auth/device/coverage.data -covermode=atomic -v
go test ./auth/addin/... -coverprofile=auth/addin/coverage.data -covermode=atomic
go test ./auth/adfs/... -coverprofile=auth/adfs/coverage.data -covermode=atomic
go test ./auth/fba/... -coverprofile=auth/fba/coverage.data -covermode=atomic
go test ./auth/ntlm/... -coverprofile=auth/ntlm/coverage.data -covermode=atomic
go test ./auth/saml/... -coverprofile=auth/saml/coverage.data -covermode=atomic
go test ./auth/tmg/... -coverprofile=auth/tmg/coverage.data -covermode=atomic
go test ./auth/anon/... -coverprofile=auth/anon/coverage.data -covermode=atomic
