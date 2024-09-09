Let's Encrypt Shell Catcher
===========================
Catches reverse shells via plain TLS or HTTPS.

Usage
-----
```
Usage: letsencryptshellcatcher [options]

Catches a TLSified reverse shell with a cert from Let's Encrypt.

Use of this program constitutes acceptance of Let's Encrypt's Terms of Service.

Options:
  -domain domain
    	TLS domain name
  -https path
    	Catch a shell via an HTTP request to the path
  -listen string
    	listen-address (default "0.0.0.0:443")
  -no-tls
    	Disable TLS, for testing
  -staging
    	Use Let's Encrypt's staging environment
```

Examples
--------

### Plain TLS
```sh
bin/letsencryptshellcatcher -domain example.com
```

### HTTPS: 
```
bin/letsencryptshellcatcher -domain example.com -https /shell
```
Make the HTTPS request to `/shell`.


