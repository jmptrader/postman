package actions

import (
	"log"

	"postman/client"
)

type AuthenticatedMsg struct {
	Id int `json:"senderId"`
}

func Authenticated(args interface{}) {
	msg := args.(*AuthenticatedMsg)
	log.Printf("client: authenticated with id %d", msg.Id)
	client.Postman.Tunnel.SetAuthenticated()
}
