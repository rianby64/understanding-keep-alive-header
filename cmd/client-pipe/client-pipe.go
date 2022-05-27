package main

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var requestBody = []byte(`GET /home HTTP/1.1
Host: localhost:8080
User-Agent: curl/7.76.1
Accept: */*

`)

func main() {
	response := make([]byte, 1000)
	client, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	closed := make(chan struct{}, 1)
	defer func() {
		closed <- struct{}{}
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()

	N := 30
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(N))

	defer cancel()

	lastBody := ""
	maxD := 0
	defer func() {
		log.Println("lastBody= ", lastBody)
		log.Println("maxD =", maxD)
	}()

	wg := &sync.WaitGroup{}
	go func() {
		for {
			n, err := client.Read(response)
			if err != nil {
				time.Sleep(time.Second)
				select {
				case <-closed:
					return
				default:
					panic(err)
				}
			}

			responseStr := string(response[:n])
			sep := "\r\n\r\n"
			w := strings.Split(responseStr, sep)
			d := len(w) - 1
			if maxD < d {
				maxD = d
			}

			body := w[len(w)-1]
			wg.Add(-d)

			lastBody = body
		}
	}()

	callIt := func() {
		if _, err := client.Write(requestBody); err != nil {
			panic(err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		N := 3
		wg.Add(N)
		for j := 0; j < N; j++ {
			go callIt()
		}

		wg.Wait()
	}

}
