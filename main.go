package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const MaxTimeOutInSecond = 2

func main() {

}

func GoCounter(urls []string) map[string]int {
	res := make(map[string]int)

	var channels []chan int
	for i := 0; i < len(urls); i++ {
		channels = append(channels, make(chan int, 1))
	}

	for i, url := range urls {
		go runParser(channels[i], url)
	}

	select {
	case <-time.After(MaxTimeOutInSecond * time.Second):
		for i, ch := range channels {
			res[urls[i]] = <-ch
		}
		return res
	}
}

func runParser(count chan int, url string) {
	client := http.Client{
		Timeout: MaxTimeOutInSecond * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		close(count)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			close(count)
			return
		}
		bodyString := string(bodyBytes)
		count <- strings.Count(bodyString, "Go")
	}
	close(count)
}

func printResult(data map[string]int) {
	sum := 0
	for key, count := range data {
		fmt.Println(key, ":", count)
		sum += count
	}
	fmt.Println(sum)
}
