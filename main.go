package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	MaxTimeOutInSecond      = 1
	CodeOfIncorrectResponse = -1
	StatusRequestTimeout    = -http.StatusRequestTimeout
	StatusBadRequest        = -http.StatusBadRequest
)

func printCountOfGo(urls []string) {
	list := GoCounter(urls)
	fmt.Println(stringResult(&list))
}

func GoCounter(urls []string) map[string]int {
	res := make(map[string]int)
	urls = uniqueList(&urls)

	chBasket := make(chan map[string]int)
	waitChannel := make(chan struct{}, runtime.NumCPU()*3)

	go func() {
		for _, url := range urls {
			waitChannel <- struct{}{}
			go func(u string) {
				runParser(chBasket, u)
				<-waitChannel
			}(url)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(urls))
	go func() {
		for ch := range chBasket {
			for k, v := range ch {
				res[k] = v
				wg.Done()
			}
		}
	}()

	wg.Wait()
	close(chBasket)

	return res
}

func runParser(count chan map[string]int, url string) {
	respond := make(map[string]int)

	client := http.Client{
		Timeout: MaxTimeOutInSecond * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			respond[url] = StatusRequestTimeout
		} else {
			respond[url] = StatusBadRequest
		}
		count <- respond
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respond[url] = CodeOfIncorrectResponse
			return
		}
		bodyString := string(bodyBytes)
		respond[url] = strings.Count(bodyString, "Go")
		count <- respond
	} else {
		respond[url] = StatusBadRequest
		count <- respond
	}
}

func stringResult(data *map[string]int) string {
	res := ""
	sum := 0
	for key, count := range *data {
		switch {
		case count >= 0:
			res += fmt.Sprintf("%s : %d \n", key, count)
			sum += count

		case count == StatusRequestTimeout:
			res += fmt.Sprintf("%s : %s", key, "StatusRequestTimeout\n")

		default:
			res += fmt.Sprintf("%s : %s", key, "StatusBadRequest\n")
		}
	}
	res += fmt.Sprintf("count : %d \n", sum)
	return res
}

func uniqueList(list *[]string) []string {
	keys := make(map[string]bool)
	uniqList := []string{}
	for _, entry := range *list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqList = append(uniqList, entry)
		}
	}
	return uniqList
}
