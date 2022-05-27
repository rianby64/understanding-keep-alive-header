package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	c := &http.Client{}

	N := 120
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(N))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/home", nil)
	if err != nil {
		panic(err)
	}

	Total := int64(0)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		res, err := c.Do(req)
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			panic(err)
		}

		if Total > 0 {
			delta := i - Total
			fmt.Println("in one second:", delta)
		}

		Total = i
		time.Sleep(time.Second)
	}
}
