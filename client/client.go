package client

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Addr string
	Conn net.Conn
}

func New(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Addr: addr,
		Conn: conn,
	}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	var buf bytes.Buffer

	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("SET"), resp.StringValue(key), resp.StringValue(value)})

	_, err := io.Copy(c.Conn, &buf)
	return err
}
