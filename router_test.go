package router

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := New()
	serverResponse := "server response"
	serverStatus := 200
	assert.Panics(t, func() {
		router.Handle("", "/", func(rw http.ResponseWriter, req *http.Request, _ Params) {})
	})

	router.Get("/", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	router.Get("/:a", func(rw http.ResponseWriter, req *http.Request, ps Params) {
		assert.NotEqual(t, ps["a"], "")
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	assert.Panics(t, func() {
		router.Get("/:a", func(rw http.ResponseWriter, req *http.Request, ps Params) {})
	})

	router.Get("/:a/b", func(rw http.ResponseWriter, req *http.Request, ps Params) {
		assert.NotEqual(t, ps["a"], "")
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	assert.Panics(t, func() {
		router.Get("/*a", func(rw http.ResponseWriter, req *http.Request, _ Params) {})
	})

	router.Get("/a/*b", func(rw http.ResponseWriter, req *http.Request, ps Params) {
		assert.NotEqual(t, ps["b"], "")
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	paths := []string{"/", "/a", "/a/b", "/ab/b"}
	for _, path := range paths {
		resp, err := http.Get(serverURL + path)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, resp.StatusCode, serverStatus)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
		resp.Body.Close()
	}
}

func TestIgnoreCase(t *testing.T) {
	router := New()
	serverResponse := "server response"
	serverStatus := 200
	router.Get("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/A/B")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, serverStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bodyBytes), serverResponse)
}

func TestTrailingSlashRedirect(t *testing.T) {
	t.Run("with slash", func(t *testing.T) {
		router := New()
		serverResponse := "server response"
		serverStatus := 200
		router.Get("/a/b/", func(rw http.ResponseWriter, req *http.Request, _ Params) {
			rw.WriteHeader(serverStatus)
			rw.Write([]byte(serverResponse))
		})

		server := httptest.NewServer(router)
		defer server.Close()
		serverURL := server.URL
		resp, err := http.Get(serverURL + "/a/b")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, serverStatus)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(bodyBytes), serverResponse)
	})

	t.Run("with slash", func(t *testing.T) {
		router := New()
		serverResponse := "server response"
		serverStatus := 200
		router.Get("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
			rw.WriteHeader(serverStatus)
			rw.Write([]byte(serverResponse))
		})

		server := httptest.NewServer(router)
		defer server.Close()
		serverURL := server.URL
		resp, err := http.Get(serverURL + "/a/b/")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, serverStatus)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(bodyBytes), serverResponse)
	})
}

func TestNotFound(t *testing.T) {
	router := New()
	serverResponse := "server response"
	serverStatus := 200
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, serverStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bodyBytes), serverResponse)
}

func TestSamePatternWidthDifferentMethod(t *testing.T) {
	router := New()
	serverResponse := "server response"
	serverStatus := 200
	router.Get("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	router.Post("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, serverStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bodyBytes), serverResponse)
}
