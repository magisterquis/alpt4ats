Code for the talk "A Lazy Programmer's Tips For Avoiding The SOC"
=================================================================
The code in this repository was used for the talk "A Lazy Programmer's Tips for
Avoiding The SOC," presented at BSides Belfast on 12 September 2024.

Slides: https://docs.google.com/presentation/d/1yUaalv-a_5oI9qYMUC7VCRqgkwIQFtlJcRWxI5TaAaE

Contents
--------
Each subdirectory of [`src`](./src) in this repository corresponds to one of
the tips presented in the talk.  They are as follows, listed in the order
presented:

Subdirectory                                                           | Description
-----------------------------------------------------------------------|------------
[`shello_world`](./src/shello_world)                                   | `Hello, World!`, shell-style
[`argv_from_source`](./src/argv_from_source)                           | [`shello_world`](./src/shello_world), but with baked-in arguments
[`argv_from_environment`](./src/argv_from_environment)                 | [`argv_from_source`](./src/argv_from_source) but taking config from the environment
[`argv_from_stdin`](./src/argv_from_stdin)                             | [`argv_from_environment`](./src/argv_from_environment) but taking config from standard input
[`argv_from_file_descriptor`](./src/argv_from_file_descriptor)         | [`argv_from_stdin`](./src/argv_from_stdin) but taking config from file descriptor 7
[`argv_from_argv0`](./src/argv_from_argv0)                             | [`argv_from_file_descriptor`](./src/argv_from_file_descriptor) but taking config from `argv[0]`
[`letsencryptshellcatcher`](./letsencryptshellcatcher)                 | Catches TLS/HTTPS reverse shells, using certs from [Let's Encrypt](https://letsencrypt.org)
[`comms_over_tls`](./src/comms_over_tls)                               | [`argv_from_file_descriptor`](./src/argv_from_file_descriptor) but with TLS comms
[`comms_over_https`](./src/comms_over_https)                           | [`comms_over_tls`](./src/comms_over_tls) but using HTTPS
[`comms_without_dns`](./src/comms_without_dns)                         | [`comms_over_https`](./src/comms_over_https) but without a DNS lookup
[`read_from_file_descriptor`](./src/read_from_file_descriptor)         | [`comms_without_dns`](./src/comms_without_dns) but reads from a file descriptor
[`map_and_regex_file_descriptor`](./src/map_and_regex_file_descriptor) | [`read_from_file_descriptor`](./src/read_from_file_descriptor) but maps the file from the file descriptor into memory and exfils regex matches
[`change_own_timestamp`](./src/change_own_timestamp)                   | [`map_and_regex_file_descriptor`](./src/map_and_regex_file_descriptor) but sets the timestamp of `argv[0]` back a week
[`injecty_lib`](./src/injecty_lib)                                     | Simpler shell than [`shello_world`](./shello_world), but injectable
[`fake_edr`](./src/fake_edr)                                           | A process which exists only to be an injection target

Building
--------
On OpenBSD, `make` should do the trick.

On other platforms, BSD make *might* work.  Failing that,
[`build.sh`](./build.sh) should be sufficient to build the binaries.

In either case, `LINKFLAGS` may be set to pass variables to the Go compiler
via `-linkflags`.  Due to shell quoting issues, single-quotesgT

### Useful Makefile Targets
Target      | Description
------------|------------
`all`       | Build ALL the things
`diffs`     | Assume sources are correct, rebuild diffs
`fromdiffs` | Assume diffs are correct, rebuild sources
`clean`     | Delete everything which can't be rebuilt
`cleanbins` | Delete binaries, leave sources
`test`      | Runs testish things, but there's no proper tests
