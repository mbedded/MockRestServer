# MockRestServer

This server can be used to store and receive any (text) data with REST


## How to compile

```
# Compile for Linux
env GOOS=linux GOARCH=amd64 go build -o test_linux

# Compile for Windows
env GOOS=windows GOARCH=amd64 go build -o test_windows

# Compile with version:
go build -ldflags "-X main.version=1.0" -o test_version
```

Build for multiple environments:
https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

Parameters/Arguments:
- port
- database
- version