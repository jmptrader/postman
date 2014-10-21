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
)

type DomainProcessor struct {
	Domain     string
	queueList  string
	sendingSet string
	available  bool
	isNew      bool
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

func CreateProcessor(domain string) {
	dp, err := NewProcessor(domain)
	if err != nil {
		return
	}
	log.Printf("Processor for domain:%s create success.", domain)
	dp.StartGuard()
}

// return the pointer for the processor ur looking for
func GetProcessor(domain string) (dp *DomainProcessor, err error) {
	dp, exist := processContainer[domain]
	if !exist {
		err = errors.New("No process found.")
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
		err = errors.New("Process has exist.")
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
		isNew:      true,
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
	dp.DisAvailable()
	return true
}

func (dp *DomainProcessor) SetAvailable() {
	dp.available = true
}

func (dp *DomainProcessor) DisAvailable() {
	dp.available = false
}

// create sender gorutine for mail
func (dp *DomainProcessor) CreateSender(m *mail.Mail) {
	if !dp.isNew {
		wait := time.NewTimer(GetDeliverInterval(dp.Domain))
		<-wait.C
	} else {
		dp.isNew = false
	}
	if m.Deliver() == nil {
		dp.SetAvailable()
	}
}

// Make a certain mail wait 30s and resend.
func (dp *DomainProcessor) delayDeliver(messageId string) {
	store.Add(dp.sendingSet, messageId)
	wait := time.NewTimer(time.Second * 30)
	<-wait.C
	store.LPush(dp.queueList, []string{messageId})
	store.Rem(dp.sendingSet, messageId)
}

// Listening to pre sending queue.
func (dp *DomainProcessor) StartGuard() {
Loop:
	messageId, err := store.BLPOP(dp.queueList, time.Minute)
	if err != nil {
		// close progress when no sending mails
		if store.Size(dp.sendingSet) < 1 {
			log.Printf("Domain sender closed: %s", dp.Domain)
			delete(processContainer, dp.Domain)
			return
		}
		goto Loop
	}
	mail, err := mail.Get(messageId)
	if err != nil {
		log.Fatalf("Get mail %s with error: %s", messageId, err.Error())
	}
	// If current mail.Domain frequency reach limit && queue left more than three, delay resend.
	// if frequencyMonitor.IsLock(mail) && store.Size(dp.queueList) > 0 {
	// 	go dp.delayDeliver(messageId)
	// 	goto Loop
	// }
	store.Add(dp.sendingSet, messageId)
	if !dp.CheckSender() {
		ticker := time.NewTicker(GetDeliverInterval(dp.Domain) / 3)
		for _ = range ticker.C {
			if dp.CheckSender() {
				ticker.Stop()
				goto Send
			}
		}
	}
Send:
	// frequencyMonitor.Record(mail)
	log.Printf("Mail: %s for %s find sender.", messageId, mail.To)
	dp.CreateSender(mail)
	goto Loop
}
