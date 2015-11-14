package serve

import (
	"log"
	"net"
	"net/http"
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
		log.Println(err)
		return
	}

	defer l.Close()
	defer os.Remove(socket)

	go func() {
		log.Println("listening on", socket)
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
