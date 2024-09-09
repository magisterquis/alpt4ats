package main

/*
 * tls_shell.go
 * Handle TLS shells
 * By J. Stuart McMurray
 * Created 20240904
 * Last Modified 20240905
 */

import (
	"log"
	"net"
)

// ServeTLS handles shells over TLS.  It terminates the program on fatal
// error.
func ServeTLS(l net.Listener) {
	/* Get a connection for the shell. */
	c, err := l.Accept()
	if nil != err {
		log.Fatalf("Error while waiting for shell connection: %s", err)
	}
	log.Printf("Got shell connection from %s", c.RemoteAddr())
	defer c.Close()
	l.Close()

	/* Proxy stdio. */
	ProxyStdio(c, c)
}
