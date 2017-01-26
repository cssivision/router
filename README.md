# Router

Router is a minimalist HTTP request router for [Go](https://golang.org/).

# Installation
```sh
go get github.com/cssivision/router
```

# Usage

```go
package main

import (
    "net/http"
    "github.com/cssivision/router"
)

func main() {
    r := router.New()
    r.Get("/", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /"))
    })

    r.Get("/a", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a"))
    })

    // path `/a` take precedence over `/:name`
    r.Get("/:name", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a/b"))
    })

    r.Get("/a/b", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a/b"))
    })

    r.Get("/:name", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a/b"))
    })

    http.ListenAndServe(":8080", r)
}
```

# TODO

* add test
* add code comment

