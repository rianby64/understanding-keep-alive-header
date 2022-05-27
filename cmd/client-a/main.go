package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	c := &http.Client{}

	defer func() {
		c.CloseIdleConnections()
	}()

	N := 30
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(N))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/home", nil)
	if err != nil {
		panic(err)
	}

	lastBody := ""
	i := 0
	func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				i++
			}

			resp, err := c.Do(req)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return
				}

				panic(err)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			lastBody = string(data)
		}
	}()

	fmt.Println(lastBody)
	fmt.Println("LEN", i)
	fmt.Println("LEN/s", i/N)
}
