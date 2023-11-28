package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := cron.New()
	eA, _ := c.AddFunc("@every 3s", func() {
		t.Log("Every 3s")
	})
	_, _ = c.AddFunc("@every 5s", func() {
		t.Log("Every 5s")
	})
	go removeEntryID(c, eA)
	c.Run()
}

func removeEntryID(c *cron.Cron, entryID cron.EntryID) {
	time.Sleep(10 * time.Second)
	c.Remove(entryID)
	fmt.Println(len(c.Entries()))
}
