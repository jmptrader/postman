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
	if m.Immediate {
		err = m.Deliver()
		m.Destroy()
		// call webHook once if mail is sent immediately
		if err != nil {
			m.CallWebHook(map[string]string{
				"event": "dropped",
				"error": err.Error(),
			})
			return
		}
		m.CallWebHook(map[string]string{
			"event": "delivered",
		})
		return
	}
	processor.ArrangeMail(m)
	return
}
