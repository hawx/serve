package serve

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/activation"
)

func Serve(port, socket string, handler http.Handler) {
	Server(port, socket, &http.Server{Handler: handler})
}

func Server(port, socket string, srv *http.Server) {
	listeners, err := activation.Listeners()
	if err != nil {
		panic(err)
	}

	if len(listeners) == 1 {
		onSocket(listeners[0], "", srv)
	} else if socket == "" {
		onPort(port, srv)
	} else {
		l, err := net.Listen("unix", socket)
		if err != nil {
			log.Println(err)
			return
		}

		onSocket(l, socket, srv)
	}
}

func onPort(port string, srv *http.Server) {
	go func() {
		srv.Addr = ":" + port
		log.Println("listening on port :" + port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	catchInterrupt(srv)
}

func onSocket(l net.Listener, socket string, srv *http.Server) {
	defer l.Close()
	if socket != "" {
		defer os.Remove(socket)
	}

	go func() {
		if socket == "" {
			log.Println("listening on systemd provided socket")
		} else {
			log.Println("listening on", socket)
		}
		if err := srv.Serve(l); err != nil {
			log.Println(err)
		}
	}()

	catchInterrupt(srv)
}

func catchInterrupt(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	s := <-c
	log.Printf("caught %s: shutting down\n", s)

	srv.Shutdown(context.Background())
}
