package service

import (
	"flag"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/golang/glog"
	"log"
	"testing"
	"time"
)

func Test_mailReader(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()
	glog.V(4).Info("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.qq.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	glog.V(4).Info("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("skeyic@foxmail.com", "uaenqonjnhnpbidf"); err != nil {
		log.Fatal(err)
	}
	glog.V(4).Info("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}
	glog.V(4).Infof("Inbox: %+v", mbox)

	criteria := imap.NewSearchCriteria()
	timeSince := "2021-02-03 00:00:00"
	tpl := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	ts, _ := time.ParseInLocation(tpl, timeSince, loc)
	glog.V(4).Infof("Search the mails since %s", ts)
	criteria.Since = ts
	//criteria.Header.Set("FROM", "ark@ark-funds.com")
	//criteria.Body = []string{"ARK Investment Management Trading Information"}

	glog.V(4).Infof("criteria: %+v", criteria.Format())

	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal("Search:", err)
	}

	glog.V(4).Infof("IDS: %v", ids)

	//// Get the last 4 messages
	//from := uint32(1)
	//to := mbox.Messages
	//if mbox.Messages > 10 {
	//	// We're using unsigned integers here, only subtract if the result is > 0
	//	from = mbox.Messages - 10
	//}
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Messages from search:")
	for msg := range messages {
		glog.V(4).Infof("Subject: %s, From: %s, Date: %s ", msg.Envelope.Subject, msg.Envelope.From[0], msg.Envelope.Date)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	glog.V(4).Info("Done!")
}
