package main

import (
	"net/http"
	"net/http/httptest"

	"reflect"
	"testing"
	"time"
)

func TimeOutServer(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2000 * time.Millisecond)
}

func BadServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	return
}

func TestGoCounter(t *testing.T) {
	var links []string

	// The bunch of servers with timeout
	countServers := 40
	servs := make([]*httptest.Server, countServers)
	expected := map[string]int{}
	for i := 0; i < countServers; i++ {
		servs[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(300 * time.Millisecond)
			w.Write([]byte(" I'am first server"))
			w.Write([]byte("Two occurrences Go Go"))
			return
		}))
		defer servs[i].Close()
		links = append(links, servs[i].URL)
		// Duplicate
		links = append(links, servs[i].URL)
		expected[servs[i].URL] = 2
	}

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(700 * time.Millisecond)
		w.Write([]byte(" I'am second server"))
		w.Write([]byte("Five occurrences Go Go Go Go Go"))
		return
	}))
	defer ts2.Close()
	links = append(links, ts2.URL)

	badServer := httptest.NewServer(http.HandlerFunc(BadServer))
	defer badServer.Close()
	links = append(links, badServer.URL)

	timeOutServer := httptest.NewServer(http.HandlerFunc(TimeOutServer))
	defer timeOutServer.Close()
	links = append(links, timeOutServer.URL)

	expected[ts2.URL] = 5
	expected[badServer.URL] = StatusBadRequest
	expected[timeOutServer.URL] = StatusRequestTimeout

	result := GoCounter(links)

	eq := reflect.DeepEqual(result, expected)
	if !eq {
		t.Errorf("Expected %v instance of %v", expected, result)
	}
	printCountOfGo(links)
}
