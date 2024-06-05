package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	Addr string
	Conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Addr: addr,
		Conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	var buf bytes.Buffer

	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("SET"), resp.StringValue(key), resp.StringValue(value)})

	_, err := c.Conn.Write(buf.Bytes())
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	var buf bytes.Buffer

	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("GET"), resp.StringValue(key)})

	_, err := c.Conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := c.Conn.Read(b)

	return string(b[:n]), err
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
