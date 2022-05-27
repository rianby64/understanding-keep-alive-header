package main

import (
	"context"
	"log"
	"net"
	"strings"
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

	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()

	N := 30
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(N))

	defer cancel()

	i := 0
	lastBody := ""
	defer func() {
		log.Println(lastBody)
		log.Println("LEN", i)
		log.Println("LEN/s", i/N)
	}()

	for {

		select {
		case <-ctx.Done():
			return
		default:
			i++
		}

		if _, err := client.Write(requestBody); err != nil {
			panic(err)
		}

		n, err := client.Read(response)
		if err != nil {
			panic(err)
		}

		responseStr := string(response[:n])
		sep := "\r\n\r\n"
		w := strings.Split(responseStr, sep)
		body := w[1]

		lastBody = body
	}

}
