// Program letsencryptshellcatcher - Catches a TLSified reverse shell with a
// cert from Let's Encrypt
package main

/*
 * letsencryptshellcatcher.go
 * Catches a TLSified reverse shell with a cert from Let's Encrypt
 * By J. Stuart McMurray
 * Created 20240904
 * Last Modified 20240906
 */

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// StagingURL is the directory URL for Let's Encrypt's staging environment.
const StagingURL = "https://acme-staging-v02.api.letsencrypt.org/directory"

func main() {
	/* Command-line flags. */
	var (
		staging = flag.Bool(
			"staging",
			false,
			"Use Let's Encrypt's staging environment",
		)
		httpsPath = flag.String(
			"https",
			"",
			"Catch a shell via an HTTP request to the `path`",
		)
		lAddr = flag.String(
			"listen",
			"0.0.0.0:443",
			"listen-address",
		)
		domain = flag.String(
			"domain",
			"",
			"TLS `domain` name",
		)
		noTLS = flag.Bool(
			"no-tls",
			false,
			"Disable TLS, for testing",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

Catches a TLSified reverse shell with a cert from Let's Encrypt.

Use of this program constitutes acceptance of Let's Encrypt's Terms of Service.

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Make sure we have a domain to serve. */
	if "" == *domain && !*noTLS {
		log.Fatalf("Need a domain (-domain)")
	}

	/* Listen for a connection. */
	var l net.Listener
	if *noTLS {
		l = listenPlaintext(*lAddr)
	} else {
		l = listenTLS(*lAddr, *domain, *staging)
	}
	defer l.Close()
	log.Printf("Listening on %s", l.Addr())

	/* Handle either plain TLS or HTTPS shells. */
	if "" != *httpsPath {
		ServeHTTPS(l, *httpsPath)
	} else {
		ServeTLS(l)
	}
}

// listenPlaintext listens for TCP connections to addr.  It terminates the
// program on error.
func listenPlaintext(lAddr string) net.Listener {
	l, err := net.Listen("tcp", lAddr)
	if nil != err {
		log.Fatalf("Error listening on %s: %s", lAddr, err)
	}
	return l
}

// listenTLS listens for TLS connections to addr, provisioning a cert from
// letsEncrypt, possibly via its staging URL.  It terminates the program on
// error.
func listenTLS(lAddr, domain string, staging bool) net.Listener {
	/* Manager to manage our cert-getting. */
	mgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}
	if staging { /* Staging URL, no cache. */
		mgr.Client = &acme.Client{DirectoryURL: StagingURL}
	} else { /* Cache directory. */
		/* Make sure we have a cache directory. */
		cdir := cacheDir()
		if err := os.MkdirAll(cdir, 0700); nil != err {
			log.Fatalf(
				"Error making certificate cache "+
					"directory %s: %s",
				cdir,
				err,
			)
		}
		mgr.Cache = autocert.DirCache(cdir)
	}
	/* Actually listen. */
	l, err := tls.Listen("tcp", lAddr, mgr.TLSConfig())
	if nil != err {
		log.Fatalf("Error listening on %s: %s", lAddr, err)
	}
	return l
}
