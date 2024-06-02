package main

import (
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

		var cmd Commander
		if v.Type() == resp.Array {
			rawCMD := v.Array()[0]
			switch rawCMD.String() {
			case CommandGET:
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
				}
			case CommandSET:
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}
			}
			p.msgChan <- Message{
				cmd:  cmd,
				peer: p,
			}
		}
	}
	return nil
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}
