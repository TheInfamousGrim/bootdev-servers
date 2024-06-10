package main

import (
	"fmt"
	"time"
)

func handleRequests(reqs <-chan request) {
	for req := range reqs {
		go handleRequest(req)
	}
}

// don't touch below this line

type request struct {
	path string
}

func main() {
	reqs := make(chan request, 100)
	go handleRequests(reqs)
	for i := 0; i < 4; i++ {
		reqs <- request{path: fmt.Sprintf("/path/%d", i)}
		time.Sleep(25 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("1 second passed, killing server")
}

func handleRequest(req request) {
	fmt.Println("Handling request for", req.path)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Done with request for", req.path)
}
