package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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

func main() {

}

func printCountOfGo(urls []string) {
	list := GoCounter(urls)
	fmt.Println(stringResult(&list))
}

func GoCounter(urls []string) map[string]int {
	res := make(map[string]int)
	urls = uniqueList(&urls)

	var channels []chan int
	for i := 0; i < len(urls); i++ {
		channels = append(channels, make(chan int, 1))
	}
	var wg sync.WaitGroup
	for i, url := range urls {
		wg.Add(1)
		go func(ch chan int, ur string) {
			defer wg.Done()
			runParser(ch, ur)
		}(channels[i], url)
	}
	wg.Wait()

	for i, ch := range channels {
		res[urls[i]] = <-ch
	}
	return res
}

func runParser(count chan int, url string) {
	client := http.Client{
		Timeout: MaxTimeOutInSecond * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			count <- StatusRequestTimeout
		} else {
			count <- StatusBadRequest
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			count <- CodeOfIncorrectResponse
			return
		}
		bodyString := string(bodyBytes)
		count <- strings.Count(bodyString, "Go")
	} else {
		count <- StatusBadRequest
	}
	close(count)
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
