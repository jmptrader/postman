package actions

import (
	"log"

	"postman/tunnel"
)

type AuthenticatedMsg struct {
	Id int `json:"senderId"`
}

func Authenticated(c *tunnel.Client, args interface{}) {
	msg := args.(*AuthenticatedMsg)
	log.Printf("client: authenticated with id %d", msg.Id)
	c.SetAuthenticated()
}
