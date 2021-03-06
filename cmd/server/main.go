package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	shutdownTimeout = time.Second * 10
)

func main() {
	server := createServer()
	log.Println("starting...")

	go mustListenAndServe(server)

	log.Println("serving!", server.Addr)

	waitShutdown(server)

	log.Println("... bye")
}

func createServer() *http.Server {
	r := http.NewServeMux()

	i := 1

	// curl -v --cacert ./cert.pem -X GET https://localhost:8080/home
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		response := []byte(fmt.Sprintf("%d", i))
		if _, err := w.Write(response); err != nil {
			log.Println(err)
		}

		i = i + 1
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func mustListenAndServe(server *http.Server) {
	listener, err := net.Listen("tcp6", server.Addr)
	if err != nil {
		log.Panicln(err)
	}

	// $ go run /usr/lib/golang/src/crypto/tls/generate_cert.go --host localhost
	if err := server.ServeTLS(listener, "./cert.pem", "./key.pem"); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panicln(err)
		}
	}
}

func waitShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err.Error(), "err := server.Close()")
	}
}
