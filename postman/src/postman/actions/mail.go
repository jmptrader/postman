package actions

import (
	"log"

	"postman/mail"
)

func SendMail(args interface{}) {
	m := args.(*mail.Mail)
	log.Printf("mail: receive new mail for %s", m.To)
	err := m.Create()
	if err != nil {
		log.Printf("mail: create mail %s", err.Error())
		return
	}
	return
}
