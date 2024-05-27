package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Addr string
}

func New(addr string) *Client {
	return &Client{
		Addr: addr,
	}
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("SET"), resp.StringValue(key), resp.StringValue(value)})
	_, err = conn.Write([]byte(buf.String()))
	return err
}
