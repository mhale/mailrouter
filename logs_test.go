package main

import (
	"fmt"
	"net"
	"testing"
)

func TestLogListAdd(t *testing.T) {
	logs := LogList{}
	ip := net.ParseIP("127.0.0.1")
	to := []string{"To"}

	// Test that log list grows to MaxLogs in size.
	for i := 1; i <= MaxLogs; i++ {
		logs.Add(ip, "From", to, fmt.Sprintf("%d", i), "Filter", "Route", "Status", "Error")
		if len(logs.Logs) != i {
			t.Errorf("LogList contains %v entries, want %v", len(logs.Logs), i)
		}
	}

	// Test that the most recently added log is in the first element.
	if logs.Logs[0].Subject != fmt.Sprintf("%d", MaxLogs) {
		t.Errorf("LogList newest entry subject is %v, want %v", logs.Logs[0].Subject, MaxLogs)
	}

	// Test that log list grows no further than MaxLogs in size.
	for i := 1; i < MaxLogs; i++ {
		logs.Add(ip, "From", to, fmt.Sprintf("%d", i), "Filter", "Route", "Status", "Error")
		if len(logs.Logs) != MaxLogs {
			t.Errorf("LogList contains %v entries, want %v", len(logs.Logs), MaxLogs)
		}
	}
}
