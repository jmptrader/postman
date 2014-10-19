package actions

import (
	"log"

	"postman/models/mail"
	"postman/tunnel"
)

func SendMail(c *tunnel.Client, args interface{}) {
	m := args.(*mail.Mail)
	log.Printf("mail: receive new mail for %s", m.To)
	err := m.Create(c)
	if err != nil {
		log.Printf("mail: create mail %s", err.Error())
		return
	}
	return
}
