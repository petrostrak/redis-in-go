package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClients(t *testing.T) {
	nClients := 10

	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}

			key := fmt.Sprintf("client_foo_%d", it)
			val := fmt.Sprintf("client_bar_%d", it)

			if err := c.Set(context.Background(), key, val); err != nil {
				log.Fatal(err)
			}

			value, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %s got this val back =>", value)
		}(i)
	}

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
