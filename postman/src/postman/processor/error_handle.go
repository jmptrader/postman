package processor

import (
	"log"
	"time"

	"postman/mail"
)

var retryInterval [3]time.Duration = [3]time.Duration{
	time.Minute * 5,
	time.Minute * 10,
	time.Minute * 25,
}

func getExceptionTreatment(l string) (string, error) {
	return tunnel.RequestBlock("exception", map[string]string{"log": l})
}

func errorHandle(dp *DomainProcessor, m *mail.Mail, l string) {
	tr, err := getExceptionTreatment(l)
	if err != nil {
		log.Printf("get treatment %s", err.Error())
		return
	}
	switch tr {
	default:
		// ignore and do nothing
		go m.CallWebHook(map[string]string{
			"event": "dropped",
			"log":   l,
		})
		dp.SetAvailable()
		SetMailSent(m)
		return

	case "resendLater":
		go func() {
			// sender will be available after one minute
			<-time.After(time.Minute * 1)
			dp.SetAvailable()
			log.Printf("domain: %s sender recover.", m.Addr())
		}()
		if m.Retries > 2 {
			go m.CallWebHook(map[string]string{
				"event": "dropped",
				"error": m.Log,
			})
			SetMailSent(m)
			return
		}
		// Set a timer to resend mail
		<-time.After(retryInterval[m.Retries])
		m.Retries += 1
		m.Update()
		log.Printf("mail: %s will resend soon.", m.Id)
		ArrangeMail(m)

	case "resendNow":
		if m.Retries > 2 {
			go m.CallWebHook(map[string]string{
				"event": "dropped",
				"error": m.Log,
			})
			SetMailSent(m)
			return
		}
		log.Printf("mail %s is resending now.", m.Id)
		m.Retries += 1
		m.Update()
		HandleMail(m)
		dp.SetAvailable()
	}
}
