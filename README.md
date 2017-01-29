# Introduction
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
        w.Write([]byte("path /\n"))
    })

    r.Get("/a", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a\n"))
    })

    // path `/a` take precedence over `/:name`
    r.Get("/a/:name", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path: /a/b, " + "name: " + ps["name"] + "\n"))
    })

    r.Get("/a/b", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path /a/b\n"))
    })

    r.Get("/file/*filepath", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path: /a/b, " + "filepath: " + ps["filepath"] + "\n"))
    })

    http.ListenAndServe(":8080", r)
}
```
