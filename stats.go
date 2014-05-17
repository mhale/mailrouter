package main

import (
	"sync"
)

type Stats struct {
	sync.RWMutex
	MsgsSent    int
	MsgsDropped int
	DataSent    int
	DataDropped int
}

func (s *Stats) Sent(size int) {
	s.Lock()
	defer s.Unlock()
	s.MsgsSent++
	s.DataSent += size
}

func (s *Stats) Dropped(size int) {
	s.Lock()
	defer s.Unlock()
	s.MsgsDropped++
	s.DataDropped += size
}
