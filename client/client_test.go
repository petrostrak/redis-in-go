package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	var wg sync.WaitGroup

	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", it)
			val := fmt.Sprintf("client_bar_%d", it)

			if err := c.Set(context.Background(), key, val); err != nil {
				log.Fatal(err)
			}

			value, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got this val back => %s\n", it, value)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func testnewclient(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		fmt.Println("set =>", fmt.Sprintf("first_name_%d", i))
		if err := c.Set(context.Background(), fmt.Sprintf("first_name_%d", i), fmt.Sprintf("last_name_%d", i)); err != nil {
			log.Fatal(err)
		}

		value, err := c.Get(context.Background(), fmt.Sprintf("first_name_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", value)
	}
}
