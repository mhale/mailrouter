package main

import (
	"net"
	"strings"
	"sync"
	"time"
)

const MaxLogs = 20

type Log struct {
	Received string
	From     string
	To       string
	Subject  string
	Filter   string
	Route    string
}

type LogList struct {
	sync.RWMutex
	Logs []Log
}

func (ll *LogList) Add(origin net.IP, from string, to []string, subject string, filter string, route string) {
	ll.Lock()
	defer ll.Unlock()

	l := Log{
		Received: time.Now().Format("2006-01-02 15:04:05"),
		From:     from,
		To:       strings.Join(to, ", "),
		Subject:  subject,
		Filter:   filter,
		Route:    route,
	}

	// Expand the log list out to the max size.
	if len(ll.Logs) < MaxLogs {
		ll.Logs = append(ll.Logs, Log{})
	}

	// Shuffle the existing log entries down one slot and put the new one in the first slot.
	copy(ll.Logs[1:], ll.Logs[0:])
	ll.Logs[0] = l
}
