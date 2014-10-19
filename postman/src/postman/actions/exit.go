package actions

import (
	"log"

	"postman/client"
)

type ExitMsg struct {
	Msg string `json:"msg"`
}

func Exit(args interface{}) {
	defer client.Postman.Tunnel.Close()
	msg := args.(*ExitMsg)
	log.Printf("client exit: %s", msg.Msg)
}
