package main

/*
 * https_shell.go
 * Handle HTTPS shells
 * By J. Stuart McMurray
 * Created 20240904
 * Last Modified 20240905
 */

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
)

// ServeHTTPS handles shells over HTTPS.  It terminates the program on fatal
// error.
func ServeHTTPS(l net.Listener, sPath string) {
	/* HTTP Server serving our one path. */
	gotShell := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc(sPath, func(w http.ResponseWriter, r *http.Request) {
		close(gotShell)
		log.Printf("Got shell request from %s", r.RemoteAddr)
		/* Try to disable buffering. */
		w.WriteHeader(http.StatusOK)
		rc := http.NewResponseController(w)
		rc.EnableFullDuplex()
		rc.Flush()

		ProxyStdio(w, r.Body)
	})
	svr := http.Server{Handler: mux}
	/* Serve until we're shut down after getting a shell. */
	go func() {
		if err := svr.Serve(l); nil != err &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(
				"Error listening for HTTPS requests: %s",
				err,
			)
		}
	}()
	/* Wait until we've got a shell, then shutdown the listener. */
	<-gotShell
	if err := svr.Shutdown(context.Background()); nil != err {
		log.Fatalf("Error shutting down HTTPS server: %s", err)
	}
}
