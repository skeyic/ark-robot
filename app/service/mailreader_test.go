package service

import (
	"flag"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/golang/glog"
	"log"
	"strings"
	"testing"
	"time"
)

func Test_mailReader(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()
	glog.V(4).Info("Connecting to server...")

	var (
		mbox *imap.MailboxStatus
		err  error
	)

	// Connect to server
	c, err := client.DialTLS("imap.qq.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	glog.V(4).Info("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	//if err := c.Login("skeyic@foxmail.com", "uaenqonjnhnpbidf"); err != nil {
	//	log.Fatal(err)
	//}
	if err := c.Login("skeyic@foxmail.com", "uaenqonjnhnpbidf"); err != nil {
		log.Fatal(err)
	}
	glog.V(4).Info("Logged in")

	//// List mailboxes
	//mailboxes := make(chan *imap.MailboxInfo, 20)
	//done := make(chan error, 1)
	//go func() {
	//	done <- c.List("", "*", mailboxes)
	//}()
	//
	//log.Println("Mailboxes:")
	//for m := range mailboxes {
	//	log.Println("* " + m.Name)
	//}
	//
	//if err := <-done; err != nil {
	//	log.Fatal(err)
	//}

	for {
		// Select INBOX
		mbox, err = c.Select("INBOX", false)
		if err != nil {
			if strings.Contains(err.Error(), "busy") {
				glog.Error("Server busy, retry in 5 seconds")
				time.Sleep(5 * time.Second)
				continue
			}
			log.Fatal(err)
		}
		glog.V(4).Infof("Inbox: %+v", mbox)
		break
	}

	//criteria := imap.NewSearchCriteria()
	//timeSince := "2021-02-04 00:00:00"
	//timeTo := "2021-02-05 00:00:00"
	//tpl := "2006-01-02 15:04:05"
	//loc, _ := time.LoadLocation("Local")
	//ts, _ := time.ParseInLocation(tpl, timeSince, loc)
	//tt, _ := time.ParseInLocation(tpl, timeTo, loc)
	//glog.V(4).Infof("Search the mails since %s to %s", ts, tt)
	//criteria.SentSince = ts
	////criteria.SentBefore = tt
	////criteria.Header.Set("FROM", "ark@ark-funds.com")
	////criteria.Body = []string{"ARK Investment Management Trading Information"}
	//
	//glog.V(4).Infof("criteria: %+v", criteria.Format())
	//
	//ids, err := c.Search(criteria)
	//if err != nil {
	//	log.Fatal("Search:", err)
	//}
	//
	//glog.V(4).Infof("IDS: %v", ids)

	// Get the last 100 messages
	var (
		messageNum uint32 = 10
		from       uint32 = 1
	)
	to := mbox.Messages
	if mbox.Messages > messageNum {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - messageNum
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, messageNum)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Messages from Fetch:")
	for msg := range messages {
		glog.V(4).Infof("Subject: %s, From: %s, Date: %s ", msg.Envelope.Subject, msg.Envelope.From[0], msg.Envelope.Date)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	glog.V(4).Info("Done!")
}
