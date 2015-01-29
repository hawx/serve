package serve

import (
	"net"
	"net/http"
	"log"
	"os"
	"os/signal"
)

func Serve(port, socket string, handler http.Handler) {
	if socket == "" {
		Port(port, handler)
	} else {
		Socket(socket, handler)
	}
}

func Socket(socket string, handler http.Handler) {
	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	defer os.Remove(socket)

	go func() {
		log.Println("listening on", socket)
		log.Fatal(http.Serve(l, handler))
	}()

	catchInterrupt()
}

func Port(port string, handler http.Handler) {
	go func() {
		log.Println("listening on port :" + port)
		log.Fatal(http.ListenAndServe(":"+port, handler))
	}()

	catchInterrupt()
}

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Printf("caught %s: shutting down", s)
}
