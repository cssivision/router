package router

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrefix(t *testing.T) {
	router := New()
	serverResponse := "server response"
	serverStatus := 200

	v1 := router.Prefix("/api/v1")
	v1.Get("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	v2 := router.Prefix("/api/v2")
	v2.Get("/a/b", func(rw http.ResponseWriter, req *http.Request, _ Params) {
		rw.WriteHeader(serverStatus)
		rw.Write([]byte(serverResponse))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/api/v1/a/b")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, serverStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bodyBytes), serverResponse)
	resp.Body.Close()
	resp, err = http.Get(serverURL + "/api/v2/a/b")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, serverStatus)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bodyBytes), serverResponse)
	resp.Body.Close()
}
