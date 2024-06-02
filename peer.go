package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn    net.Conn
	msgChan chan Message
}

func NewPeer(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

func (p *Peer) read() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, val := range v.Array() {
				switch val.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid get command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
					}
					fmt.Printf("GET cmd: %+v\n", cmd)
				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid set command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					fmt.Printf("SET cmd: %+v\n", cmd)
				}
			}
		}
	}
	return nil
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}
