package main

import (
	"sync"
)

type Stats struct {
	sync.RWMutex
	MsgsSent    int
	MsgsDropped int
	MsgsFailed  int
	DataSent    int
	DataDropped int
	DataFailed  int
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

func (s *Stats) Failed(size int) {
	s.Lock()
	defer s.Unlock()
	s.MsgsFailed++
	s.DataFailed += size
}
