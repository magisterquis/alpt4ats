Argv from Source
================
Nearly the same as [`shello_world`](../shello_world) but instead of
command-line arguments variables baked into the source are used.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                                    | Default          | Description
--------------------------------------------|------------------|------------
[`main.Address`](./argv_from_source.go#L21) | `localhost:4444` | TCP Server Address
[`main.File`](./argv_from_source.go#L22)    | `/etc/hosts`     | File to send on the TCP connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```
