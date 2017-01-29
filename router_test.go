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

	// router.Get("/a/*b", func(rw http.ResponseWriter, req *http.Request, ps Params) {
	//     assert.NotEqual(t, ps["a"], "")
	//     rw.WriteHeader(serverStatus)
	//     rw.Write([]byte(serverResponse))
	// })

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	paths := []string{"/", "/a"}
	for _, path := range paths {
		resp, err := http.Get(serverURL + path)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, resp.StatusCode, serverStatus)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
		resp.Body.Close()
	}
}
