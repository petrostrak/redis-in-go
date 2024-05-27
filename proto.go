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
)

type Commander interface{}

type SetCommand struct {
	key, val string
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
		fmt.Printf("Read %s\n", v.Type())

		var cmd Commander
		if v.Type() == resp.Array {
			for _, val := range v.Array() {
				switch val.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid set command")
					}
					cmd = SetCommand{
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}
					return cmd, nil
				default:
				}
				fmt.Printf("%v\n", v)
			}
		}
	}
	return SetCommand{}, nil
}
