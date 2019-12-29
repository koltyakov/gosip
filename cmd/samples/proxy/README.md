# SharePoint Proxy for dev toolchains

## Start

```bash
go run ./cmd/proxy -strategy adfs -config ./config/private.json -port 9090
```

## HTTPS

```bash
openssl genrsa -out ./config/certs/private.key 2048
openssl req -new -x509 -sha256 -key ./config/certs/private.key -out ./config/certs/public.crt -days 3650
```

```bash
go run ./cmd/proxy -strategy adfs -config ./config/private.onprem-wap-adfs.json -port 443 -sslKey ./config/certs/private.key -sslCert ./config/certs/public.crt
```