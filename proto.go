package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
)

type Commander interface{}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

func parseCommand(msg string) (Commander, error) {
	rd := resp.NewReader(bytes.NewBufferString(msg))
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
						return nil, fmt.Errorf("invalid get command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
					}
					return cmd, nil
				case CommandSET:
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid set command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					return cmd, nil
				}
			}
		}
	}
	return SetCommand{}, fmt.Errorf("invalid or unknown command: %s\n", msg)
}
