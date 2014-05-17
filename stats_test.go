package main

import (
	"math/rand"
	"testing"
	"time"
)

// Record the sending of two randomly sized emails and verify that the
// sent stats show 2 emails and the combined data size.
func TestStatsSent(t *testing.T) {
	var r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var n1 = r.Intn(100)
	var n2 = r.Intn(100)

	stats := Stats{}
	stats.Sent(n1)
	if stats.MsgsSent != 1 {
		t.Errorf("stats.MsgsSent = %v, want %v", stats.MsgsSent, 1)
	}
	if stats.MsgsDropped != 0 {
		t.Errorf("stats.MsgsDropped = %v, want %v", stats.MsgsDropped, 0)
	}
	if stats.DataSent != n1 {
		t.Errorf("stats.DataSent = %v, want %v", stats.DataSent, n1)
	}
	if stats.DataDropped != 0 {
		t.Errorf("stats.DataDropped = %v, want %v", stats.DataDropped, 0)
	}

	stats.Sent(n2)
	if stats.MsgsSent != 2 {
		t.Errorf("stats.MsgsSent = %v, want %v", stats.MsgsSent, 2)
	}
	if stats.MsgsDropped != 0 {
		t.Errorf("stats.MsgsDropped = %v, want %v", stats.MsgsDropped, 0)
	}
	if stats.DataSent != n1+n2 {
		t.Errorf("stats.DataSent = %v, want %v", stats.DataSent, n1+n2)
	}
	if stats.DataDropped != 0 {
		t.Errorf("stats.DataDropped = %v, want %v", stats.DataDropped, 0)
	}
}

// Record the dropping of two randomly sized emails and verify that the
// dropped stats show 2 emails and the combined data size.
func TestStatsDropped(t *testing.T) {
	var r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var n1 = r.Intn(100)
	var n2 = r.Intn(100)

	stats := Stats{}
	stats.Dropped(n1)
	if stats.MsgsSent != 0 {
		t.Errorf("stats.MsgsSent = %v, want %v", stats.MsgsSent, 0)
	}
	if stats.MsgsDropped != 1 {
		t.Errorf("stats.MsgsDropped = %v, want %v", stats.MsgsDropped, 1)
	}
	if stats.DataSent != 0 {
		t.Errorf("stats.DataSent = %v, want %v", stats.DataSent, 0)
	}
	if stats.DataDropped != n1 {
		t.Errorf("stats.DataDropped = %v, want %v", stats.DataDropped, n1)
	}

	stats.Dropped(n2)
	if stats.MsgsSent != 0 {
		t.Errorf("stats.MsgsSent = %v, want %v", stats.MsgsSent, 0)
	}
	if stats.MsgsDropped != 2 {
		t.Errorf("stats.MsgsDropped = %v, want %v", stats.MsgsDropped, 2)
	}
	if stats.DataSent != 0 {
		t.Errorf("stats.DataSent = %v, want %v", stats.DataSent, 0)
	}
	if stats.DataDropped != n1+n2 {
		t.Errorf("stats.DataDropped = %v, want %v", stats.DataDropped, n1+n2)
	}
}
