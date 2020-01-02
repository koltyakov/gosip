# Tests

## Run automated tests

Create auth credentials store files in `./config` folder for corresponding strategies:

- private.onprem-adfs.json
- private.onprem-fba.json
- private.onprem-ntlm.json
- private.onprem-tmg.json
- private.onprem-wap-adfs.json
- private.onprem-wap.json
- private.spo-addin.json
- private.spo-user.json
- private.spo-adfs.json

Auth configs should have the same structure as [node-sp-auth's](https://github.com/s-kainet/node-sp-auth) configs. See [samples](./config/samples).

```bash
go test ./... -v -race -count=1
```

Not provided auth configs are ignored and not skipped in tests.

## Run manual test

Modify `cmd/gosip/main.go` to include required scenarios and run:

```bash
go run cmd/gosip/main.go
```

## Run CI tests

Configure environment variables:

- SPAUTH_SITEURL
- SPAUTH_CLIENTID
- SPAUTH_CLIENTSECRET
- SPAUTH_USERNAME
- SPAUTH_PASSWORD

```bash
go test ./... -race -timeout 30s
```