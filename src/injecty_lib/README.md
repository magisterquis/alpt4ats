Injecty Library
===============
Nearly the same as [`argv_from_source`](../argv_from_source) but meant for
injecting into another process.

Also works as a standalone binary.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                               | Default          | Description
---------------------------------------|------------------|------------
[`main.Address`](./injecty_lib.go#L21) | `localhost:4444` | TLS Server Address:Port
[`main.File`](./injecty_lib.go#L22)    | `/etc/hosts`     | File to send on the TLS connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:443" make
```

Building
--------
Building as a standalone binary is just like any other Go program except that
the server address will need to be set at compile-time.

Something like this works:
```sh
go build -ldflags '-X main.Address=example.com:443'
```

Building as an injectable library is not much different, but requires a C
compiler to be present (`apt -y install build-essential`, on Debian), and to
ask Go to make a shared object file with `-buildmode c-shared`, and looks like
```sh
go build -v -buildmode=c-shared -ldflags '-X main.Address=example.com:443'
```
