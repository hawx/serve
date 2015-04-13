# servers-gonna-serve [![docs](http://godoc.org/github.com/hawx/serve?status.svg)](http://godoc.org/github.com/hawx/serve)

A small package that wraps up serving a `http.Handler` via a port or socket. It
catches interrupts so any `defer`s will run properly.

``` go
package main

import (
  "github.com/hawx/serve"
  "net/http"
  "flag"
)

var (
  port   = flag.String("port", "8080", "")
  socket = flag.String("socket", "", "")
)

func main() {
  flag.Parse()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // ...
  })

  serve.Serve(*port, *socket, http.DefaultServeMux)
}
```
