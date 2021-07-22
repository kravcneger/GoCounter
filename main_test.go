package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func ValidServer(bodyResponce string, timeOut int) {
	time.Sleep(time.Duration(timeOut) * time.Millisecond)
}

func TimeOutServer(w http.ResponseWriter, r *http.Request) {
	time.Sleep(20 * time.Second)
}

func BadServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	return
}

func TestGoCounter(t *testing.T) {
	var servers []*httptest.Server
	response := "some text Go"
	for i := 1; i <= 5; i++ {
		//Каждый последующий сервер будет отвечать с задержкой N*200 милисекунд и отдавать
		//Строчку увеличенную на "some text Go"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(i*200) * time.Millisecond)
			w.Write([]byte(" I'am " + strconv.Itoa(i) + " server"))
			w.Write([]byte(response))
			return
		}))
		defer ts.Close()
		response += response
		servers = append(servers, ts)
	}

	badServer := httptest.NewServer(http.HandlerFunc(BadServer))
	defer badServer.Close()

	timeOutServer := httptest.NewServer(http.HandlerFunc(TimeOutServer))
	defer timeOutServer.Close()

	var links []string
	expected := make(map[string]int)

	for i, serv := range servers {
		links = append(links, serv.URL)
		expected[serv.URL] = i
	}
	links = append(links, badServer.URL)
	links = append(links, timeOutServer.URL)

	result := GoCounter(links)

	eq := reflect.DeepEqual(result, expected)
	if !eq {
		t.Errorf("Expected %v instance of %v", expected, result)
	}

}
