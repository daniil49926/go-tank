package main

import (
	"fmt"
	"net/http"
	"time"
)

type Settings struct {
	UrlPath  string
	Duration int
}

func knockingOnTheServer(s Settings) {
	to := time.After(1 * time.Second)
	resChan := make(chan string)
	stopChannel := make(chan bool, 1)

	go func() {
		defer fmt.Println("Stop")
		defer close(resChan)
		for {
			select {
			case <-to:
				fmt.Println("Time is up")
				stopChannel <- true
				time.Sleep(10 * time.Second)
				return
			default:
				go nonBlockingGet(s.UrlPath, resChan)
			}
		}
	}()
	countElem := 0
	for _ = range resChan {
		countElem += 1
	}
	<-stopChannel

	fmt.Println(countElem)
}

func nonBlockingGet(_url string, rc chan string) {
	resp, err := http.Get(_url)
	if err != nil {
		return
	}
	rc <- resp.Status
}
