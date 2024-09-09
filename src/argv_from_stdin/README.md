Argv from Standard Input
========================
Nearly the same as [`argv_from_environment`](../argv_from_environment) but an
attempt will be made to read configuration from stdin as well.

Standard Input
--------------
The address and file are read as whitespace-separated words from stdin, and
override the environment variables.
If only one word is read, it is taken to be the address.

EOF is necessary; if neither is to be read, redirect stdin from /dev/null.

### Examples
Set the address and file via stdin:
```sh
echo localhost:5555 /etc/passwd | ./argv_from_stdin 
```
Set just the address via stdin:
```sh
echo localhost:5555 | ./argv_from_stdin 
```
Set neither via stdin:
```sh
./argv_from_stdin </dev/null
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
[`main.Address`](./argv_from_stdin.go#L26)       | `localhost:4444`   | TCP Server Address
[`main.AddressEnvVar`](./argv_from_stdin.go#L27) | `ALPT4ATS_ADDRESS` | Environment variable which sets the TCP Server Address
[`main.File`](./argv_from_stdin.go#L28)          | `/etc/hosts`       | File to send on the TCP connection
[`main.FileEnvVar`](./argv_from_stdin.go#L29)    | `ALPT4ATS_FILE`    | Environment variable which sets the File to send on the TCP connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```

