package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func ValidServer(bodyResponce string, timeOut int) {
	time.Sleep(time.Duration(timeOut) * time.Millisecond)
}

func TimeOutServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func BadServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	return
}

func TestGoCounter(t *testing.T) {
	var links []string

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte(" I'am first server"))
		w.Write([]byte("Two occurrences Go Go"))
		return
	}))
	defer ts1.Close()
	links = append(links, ts1.URL)

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
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

	expected := map[string]int{
		ts1.URL:           2,
		ts2.URL:           5,
		badServer.URL:     0,
		timeOutServer.URL: 0,
	}

	result := GoCounter(links)

	eq := reflect.DeepEqual(result, expected)
	if !eq {
		t.Errorf("Expected %v instance of %v", expected, result)
	}

	printResult(result)
}
