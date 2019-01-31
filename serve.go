package serve

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/coreos/go-systemd/activation"
)

func Serve(port, socket string, handler http.Handler) {
	listeners, err := activation.Listeners()
	if err != nil {
		panic(err)
	}

	if len(listeners) == 1 {
		Socket(listeners[0], handler, "")
	} else if socket == "" {
		Port(port, handler)
	} else {
		l, err := net.Listen("unix", socket)
		if err != nil {
			log.Println(err)
			return
		}

		Socket(l, handler, socket)
	}
}

func Socket(l net.Listener, handler http.Handler, socket string) {
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
		if err := http.Serve(l, handler); err != nil {
			log.Println(err)
		}
	}()

	catchInterrupt()
}

func Port(port string, handler http.Handler) {
	go func() {
		log.Println("listening on port :" + port)
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			log.Println(err)
		}
	}()

	catchInterrupt()
}

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Printf("caught %s: shutting down", s)
}
