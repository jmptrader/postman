package processor

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"postman/client"
	"postman/mail"
)

const (
	sendingPrefix = "jianxin-sending-set:" // current sending mailIds for domain
	queuePrefix   = "jianxin-queue:"       // current pending deliver mailIds queue
)

var (
	processorLock    = &sync.Mutex{}
	processContainer = make(map[string]*DomainProcessor)
	store            = client.Postman.Store
	tunnel           = client.Postman.Tunnel
)

type DomainProcessor struct {
	Domain     string
	queueList  string
	sendingSet string
	available  bool
	queryLock  *sync.Mutex
}

func init() {
	keys := store.Keys(sendingPrefix)
	// Moving sending mailIds back to queue.
	for _, sendingKey := range keys {
		domain := strings.SplitN(sendingKey, ":", 2)[1]
		sendingIds := store.Members(sendingKey)
		store.LPush(queuePrefix+domain, sendingIds)
		store.Destroy(sendingKey)
	}
	// Create processor for all queue.
	queueKeys := store.Keys(queuePrefix)
	for _, queueKey := range queueKeys {
		domain := strings.SplitN(queueKey, ":", 2)[1]
		go CreateProcessor(domain)
	}
	log.Print("Main processor init success")
}

// set mail to pre sending queue
func ArrangeMail(m *mail.Mail) {
	addr := m.Addr()
	size, err := store.RPush(queuePrefix+addr, []string{m.Id})
	if err != nil {
		log.Fatalf("processor: arrange mail %s", err.Error())
	}
	if size == 1 {
		go CreateProcessor(addr)
	}
	return
}

func HandleMail(m *mail.Mail) {
	addr := m.Addr()
	size, err := store.LPush(queuePrefix+addr, []string{m.Id})
	if err != nil {
		log.Fatalf("processor: handle mail %s", err.Error())
	}
	if size == 1 {
		go CreateProcessor(addr)
	}
	return
}

func SetMailSent(m *mail.Mail) {
	m.Destroy()
	addr := m.Addr()
	domainProcessor, err := GetProcessor(addr)
	if err != nil {
		return
	}
	store.Rem(domainProcessor.sendingSet, m.Id)
}

func CreateProcessor(domain string) {
	dp, err := NewProcessor(domain)
	if err != nil {
		return
	}
	log.Printf("processor for domain:%s create success.", domain)
	dp.StartGuard()
}

// return the pointer for the processor ur looking for
func GetProcessor(domain string) (dp *DomainProcessor, err error) {
	dp, exist := processContainer[domain]
	if !exist {
		err = errors.New("no process found.")
		return
	}
	return
}

// Create new domain|sender processor
// if exist will return error
func NewProcessor(domain string) (dp *DomainProcessor, err error) {
	processorLock.Lock()
	defer processorLock.Unlock()
	_, exist := processContainer[domain]
	if exist {
		err = errors.New("process has exist.")
		return
	}
	if err != nil {
		return
	}
	dp = &DomainProcessor{
		Domain:     domain,
		queueList:  queuePrefix + domain,
		sendingSet: sendingPrefix + domain,
		available:  true,
		queryLock:  &sync.Mutex{},
	}
	processContainer[domain] = dp
	return
}

// Check if available sender for current domain
func (dp *DomainProcessor) CheckSender() bool {
	dp.queryLock.Lock()
	defer dp.queryLock.Unlock()
	if !dp.available {
		return false
	}
	dp.available = false
	return true
}

func (dp *DomainProcessor) SetAvailable() {
	dp.queryLock.Lock()
	defer dp.queryLock.Unlock()
	dp.available = true
}

// create sender gorutine for mail
func (dp *DomainProcessor) CreateSender(m *mail.Mail) {
	err := m.Deliver()
	if err == nil {
		m.CallWebHook(map[string]string{
			"event": "delivered",
		})
		dp.SetAvailable()
		SetMailSent(m)
		return
	}
	errorHandle(dp, m, err.Error())
}

// Listening to pre sending queue.
func (dp *DomainProcessor) StartGuard() {
Loop:
	messageId, err := store.BLPOP(dp.queueList, time.Minute)
	if err != nil {
		// close progress when no sending mails
		if store.Size(dp.sendingSet) < 1 {
			log.Printf("domain sender closed: %s", dp.Domain)
			delete(processContainer, dp.Domain)
			return
		}
		goto Loop
	}
	mail, err := mail.Get(messageId)
	if err != nil {
		log.Fatalf("get mail %s with error: %s", messageId, err.Error())
	}
	store.Add(dp.sendingSet, messageId)
	if !dp.CheckSender() {
		ticker := time.NewTicker(GetDeliverInterval(dp.Domain) / 3)
		// TODO: loop here can be replaced with channel someday.
		for _ = range ticker.C {
			if dp.CheckSender() {
				ticker.Stop()
				goto Send
			}
		}
	}
Send:
	log.Printf("mail: %s for %s find sender.", messageId, mail.To)
	go dp.CreateSender(mail)
	wait := time.NewTimer(GetDeliverInterval(dp.Domain))
	log.Print(wait)
	<-wait.C
	goto Loop
}
