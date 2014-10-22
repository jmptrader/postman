package actions

import (
	"log"

	"postman/mail"
	"postman/processor"
)

func SendMail(args interface{}) {
	m := args.(*mail.Mail)
	log.Printf("mail: receive new mail for %s", m.To)
	err := m.Create()
	if err != nil {
		log.Printf("mail: create mail %s", err.Error())
		return
	}
	go func() {
		if m.Immediate {
			m.Deliver()
			processor.SetMailSent(m)
			return
		}
		processor.ArrangeMail(m)
	}()
	return
}
