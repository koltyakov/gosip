# Manual testing

## Authentication configs

Create auth credentials store files in `./config` folder for corresponding strategies.

See [samples](../../config/samples).

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
