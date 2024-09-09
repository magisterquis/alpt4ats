package main

/*
 * stdio.go
 * Proxy to/from stdio
 * By J. Stuart McMurray
 * Created 20240905
 * Last Modified 20240905
 */

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ProxyStdio proxies between the shell and stdio.
func ProxyStdio(toShell io.Writer, fromShell io.Reader) {
	/* Start proxying going. */
	ech := make(chan error, 2)
	go proxy(ech, toShell, os.Stdin, "Stdin -> Shell")
	go proxy(ech, os.Stdout, fromShell, "Shell -> Stdout")

	/* Wait for it to end. */
	if err := <-ech; nil != err {
		log.Printf("Shell finished with error: %s", err)
		return
	}
	log.Printf("Shell finished")
}

// proxy copies from src to dst, stopping when the context is done, after
// which one more read may succeed.  This is kooky, but works around not being
// able to reliably stop a read from os.Stdin.
func proxy(
	ech chan<- error,
	dst io.Writer,
	src io.Reader,
	dir string,
) {
	/* We do it manuallyish to also flush. */
	b := make([]byte, 1024)
	for {
		/* Read a chunk. */
		nr, err := src.Read(b)
		if nil != err && 0 == nr {
			/* Normal ending. */
			if errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			ech <- fmt.Errorf("read %s: %w", dir, err)
			return
		}
		/* Write a chunk. */
		if _, err := dst.Write(b[:nr]); nil != err {
			ech <- fmt.Errorf("write %s: %w", dir, err)
			return
		}
		/* Actually send it. */
		if f, ok := dst.(http.Flusher); ok {
			f.Flush()
		}
	}

	ech <- nil
}
