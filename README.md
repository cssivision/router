# Introduction
Router is a minimalist HTTP request router for [Go](https://golang.org/).

# Feature
* Named parameters
* Wildcard parameters
* Trailing slash redirect
* Case sensitive
* Prefix support

# Installation
```sh
go get github.com/cssivision/router
```

# Usage

## simple pattern
```go
package main

import (
    "net/http"
    "github.com/cssivision/router"
)

func main() {
    r := router.New()
    r.Get("/", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("index\n"))
    })

    r.Get("/a/b", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("hello world!\n"))
    })

    http.ListenAndServe(":8080", r)
}
```
## named pattern
```go
package main

import (
    "net/http"
    "github.com/cssivision/router"
)

func main() {
    r := router.New()

    r.Get("/a/:name", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path: /a/:name, " + "name: " + ps["name"] + "\n"))
    })

    http.ListenAndServe(":8080", r)
}
```

## wildcard pattern
```go
package main

import (
    "net/http"
    "github.com/cssivision/router"
)

func main() {
    r := router.New()

    r.Get("/file/*filepath", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("path: /file/*filepath, " + "filepath: " + ps["filepath"] + "\n"))
    })

    http.ListenAndServe(":8080", r)
}
```

## prefix pattern
```go
package main

import (
    "net/http"
    "github.com/cssivision/router"
)

func main() {
    r := router.New()

    v1 := r.Prefix("/api/v1")
    v1.Get("/a", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("api v1\n"))
    })

    v2 := r.Prefix("/api/v2")
    v2.Get("/a", func(w http.ResponseWriter, r *http.Request, ps router.Params){
        w.Write([]byte("api v2\n"))
    })

    http.ListenAndServe(":8080", r)
}
```

## Named parameters
Named parameters only match a single path segment:
```
Pattern: /user/:name

 /user/gordon              match
 /user/you                 match
 /user/gordon/profile      no match
 /user/                    no match

Pattern: /:user/:name

 /a/gordon                 match
 /b/you                    match
 /user/gordon/profile      no match
 /user/                    no match
```
## Wildcard parameters
Match everything, therefore they must always be at the end of the pattern:

```
Pattern: /src/*filepath

 /src/                     match
 /src/somefile.go          match
 /src/subdir/somefile.go   match
 ```
## Trailing slash redirect
* TrailingSlashRedirect: /a/b/ -> /a/b
* TrailingSlashRedirect: /a/b -> /a/b/

## Case sensitive
* `IgnoreCase = true`: /A/B/ -> /a/b
* `IgnoreCase = false`: case sensitive

# Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/router/blob/master/LICENSE).
