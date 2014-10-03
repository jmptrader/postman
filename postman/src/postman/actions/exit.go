package actions

import (
	"log"

	"postman/tunnel"
)

type ExitMsg struct {
	Msg string `json:"msg"`
}

func Exit(c *tunnel.Client, args interface{}) {
	defer c.Close()
	msg := args.(*ExitMsg)
	log.Printf("client exit: %s", msg.Msg)
}
