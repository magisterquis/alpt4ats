Argv from File Descriptor
=========================
Nearly the same as [`argv_from_stdin`](../argv_from_stdin) but an
attempt will be made to read configuration a specific file descriptor number
instead of standard input (i.e. file descriptor number 0).

The file descriptor number is hardcoded to 7.

File Descriptor
---------------
The address and file as whitespace-separated words from file descriptor 7, and
override the environment variables.  If only one word is read, it is taken to
be the address.

### Examples
Set the address and file via a heredoc:
```sh
./argv_from_file_descriptor 7<<_eof
localhost:5555 /etc/passwd
_eof
```

Set just the address via bash's Here String:
```sh
./bin/argv_from_file_descriptor 7<<<"localhost:5555"
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

Variable                                                   | Default            | Description
-----------------------------------------------------------|--------------------|------------
[`main.Address`](./argv_from_file_descriptor.go#L25)       | `localhost:4444`   | TCP Server Address
[`main.AddressEnvVar`](./argv_from_file_descriptor.go#L26) | `ALPT4ATS_ADDRESS` | Environment variable which sets the TCP Server Address
[`main.File`](./argv_from_file_descriptor.go#L27)          | `/etc/hosts`       | File to send on the TCP connection
[`main.FileEnvVar`](./argv_from_file_descriptor.go#L28)    | `ALPT4ATS_FILE`    | Environment variable which sets the File to send on the TCP connection

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=example.com:8080" make
```
