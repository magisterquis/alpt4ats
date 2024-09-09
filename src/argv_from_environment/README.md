Argv from Environment
=====================
Nearly the same as [`argv_from_source`](../argv_from_source) but environment
variables may be used to override the compile-time config.

Environment Variables
---------------------
The address and file may be set with environment variables, overriding
compile-time settings, as follows

Default Variable   | Description
-------------------|------------
`ALPT4ATS_ADDRESS` | TCP Server Address
`ALPT4ATS_FILE`    | File to send on the TCP connection

The variable names may be changed at compile-time, as below.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                                               | Default            | Description
-------------------------------------------------------|--------------------|------------
[`main.Address`](./argv_from_environment.go#L23)       | `localhost:4444`   | TCP Server Address
[`main.AddressEnvVar`](./argv_from_environment.go#L24) | `ALPT4ATS_ADDRESS` | Environment variable which sets the TCP Server Address
[`main.File`](./argv_from_environment.go#L25)          | `/etc/hosts`       | File to send on the TCP connection
[`main.FileEnvVar`](./argv_from_environment.go#L26)    | `ALPT4ATS_FILE`    | Environment variable which sets the File to send on the TCP connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```
