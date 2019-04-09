package main

import (
	"fmt"
	"flag"
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"time"
)

type responseInfo struct {
	status int
	bytes int64
	duration time.Duration
}

func main() {
	// Parse arguments
	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	flag.Parse()
	if flag.NArg() == 0 || *requests == 0 || *requests < *concurrency {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	link := flag.Arg(0)
	fmt.Println(*requests, *concurrency)

	// Create pools
	jobs := make(chan string, *requests)
	results := make(chan responseInfo)
	for j := int64(1); j <= *requests; j++ {
		jobs <- link
	}
	close(jobs)

	for i := int64(1); i <= *concurrency; i++ {
		go worker(i, jobs, results)
	}

	count := int64(0)
	for response := range results {
		fmt.Println(response)
		count++
		if (count >= *requests) {
			break
		}
	}
}

func worker(workerId int64, jobs <- chan string, results chan <- responseInfo) {
	for j := range jobs {
		fmt.Println("worker", workerId, "started job", j)
		response := checkLink(j)
		fmt.Println("worker", workerId, "finished job", j)
		results <- response
	}
}

func checkLink(link string) responseInfo {
	start := time.Now()
	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)
	return responseInfo{
		status: res.StatusCode,
		bytes: read,
		duration: time.Now().Sub(start),
	}
}