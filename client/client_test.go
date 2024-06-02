package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		fmt.Println("SET =>", fmt.Sprintf("first_name_%d", i))
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
