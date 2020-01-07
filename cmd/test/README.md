# Manual testing

## Authentication configs

Create auth credentials store files in `./config` folder for corresponding strategies:

* [private.onprem-adfs.json](../auth/strategies/adfs.md#on-premises-configuration)
* [private.onprem-fba.json](../auth/strategies/fba.md#json)
* [private.onprem-ntlm.json](../auth/strategies/ntlm.md#json)
* [private.onprem-tmg.json](../auth/strategies/tmg.md#json)
* [private.onprem-wap-adfs.json](../auth/strategies/adfs.md#on-premises-behing-wap-configuration)
* [private.onprem-wap.json](../auth/strategies/adfs.md#on-premises-behing-wap-configuration)
* [private.spo-addin.json](../auth/strategies/addin.md#json)
* [private.spo-user.json](../auth/strategies/saml.md#json)
* [private.spo-adfs.json](../auth/strategies/adfs.md#sharepoint-online-configuration)

See [samples](./config/samples).

## Run manual tests

Modify `cmd/test/main.go` to include required scenarios and run:

```bash
go run ./cmd/test
```

Optionally, you can provide a strategy to use with a corresponding flag:

```bash
go run ./cmd/test -strategy adfs
```

## See also [testing section](https://go.spflow.com/contributing/testing) in docs.