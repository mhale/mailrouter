package main

import (
	"fmt"
	"net"
	"testing"
)

// Test for valid generation of summary string.
// Only populated fields should be returned.
func TestFilterSummarise(t *testing.T) {
	f := Filter{}
	summary := ""
	if output := f.Summarise(); output != summary {
		t.Errorf("filter.Summarise() = %s, want %s", output, summary)
	}
	f.From = "sender@example.com"
	summary = fmt.Sprintf("From: %s", f.From)
	if output := f.Summarise(); output != summary {
		t.Errorf("filter.Summarise() = %s, want %s", output, summary)
	}
	f.To = "recipient@example.com"
	summary = fmt.Sprintf("From: %s, To: %s", f.From, f.To)
	if output := f.Summarise(); output != summary {
		t.Errorf("filter.Summarise() = %s, want %s", output, summary)
	}
	f.Subject = "Lorem ipsum dolor sit amet"
	summary = fmt.Sprintf("From: %s, To: %s, Subject: %s", f.From, f.To, f.Subject)
	if output := f.Summarise(); output != summary {
		t.Errorf("filter.Summarise() = %s, want %s", output, summary)
	}
	f.Origin = "127.0.0.1"
	summary = fmt.Sprintf("From: %s, To: %s, Subject: %s, Origin: %s", f.From, f.To, f.Subject, f.Origin)
	if output := f.Summarise(); output != summary {
		t.Errorf("filter.Summarise() = %s, want %s", output, summary)
	}
}

func TestFilterMatch(t *testing.T) {
	tests := []struct {
		f   Filter
		out bool
	}{
		{Filter{}, false},
		// Full field positive matches
		{Filter{From: "sender@example.com"}, true},
		{Filter{From: "sender@example.com", To: "recipient@example.com"}, true},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet"}, true},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet", Origin: "127.0.0.1"}, true},
		{Filter{To: "recipient@example.com"}, true},
		{Filter{To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet"}, true},
		{Filter{To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet", Origin: "127.0.0.1"}, true},
		{Filter{Subject: "Lorem ipsum dolor sit amet"}, true},
		{Filter{Subject: "Lorem ipsum dolor sit amet", Origin: "127.0.0.1"}, true},
		{Filter{Origin: "127.0.0.1"}, true},
		// Full field negative matches
		{Filter{From: "sender2@example.com"}, false},
		{Filter{From: "sender@example.com", To: "recipient2@example.com"}, false},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet 2"}, false},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet", Origin: "127.0.0.2"}, false},
		// Partial field positive matches
		{Filter{From: "sender", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet"}, true},
		{Filter{From: "sender@example.com", To: "recipient", Subject: "Lorem ipsum dolor sit amet"}, true},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "Lorem ipsum"}, true},
		// Partial field negative matches
		{Filter{From: "sender2", To: "recipient@example.com", Subject: "Lorem ipsum dolor sit amet"}, false},
		{Filter{From: "sender@example.com", To: "recipient2", Subject: "Lorem ipsum dolor sit amet"}, false},
		{Filter{From: "sender@example.com", To: "recipient@example.com", Subject: "lorem ipsum"}, false},
	}
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	subject := "Lorem ipsum dolor sit amet"
	originIP := net.ParseIP("127.0.0.1")
	for _, tt := range tests {
		if x := tt.f.Match(from, to, subject, originIP); x != tt.out {
			t.Errorf("Filter{%v}.Match(%v, %v, %v, %v) = %v, want %v", tt.f, from, to, subject, originIP, x, tt.out)
		}
	}
}

func TestFilterMatchFrom(t *testing.T) {
	tests := []struct {
		from string
		out  bool
	}{
		{"", false},
		{"sender", true},
		{"example.com", true},
		{"example.org", false},
		{"sender@example.com", true},
	}
	f := Filter{}
	from := "sender@example.com"
	for _, tt := range tests {
		f.From = tt.from
		if x := f.MatchFrom(from); x != tt.out {
			t.Errorf("Filter{From: %s}.MatchFrom(%s) = %v, want %v", tt.from, from, x, tt.out)
		}
	}
}

func TestFilterMatchTo(t *testing.T) {
	tests := []struct {
		to  string
		out bool
	}{
		{"", false},
		{"recipient", true},
		{"recipient3", false},
		{"example.com", true},
		{"example.org", false},
		{"recipient@example.com", true},
		{"recipient2@example.net", true},
	}
	f := Filter{}
	to := []string{"recipient@example.com", "recipient2@example.net"}
	for _, tt := range tests {
		f.To = tt.to
		if x := f.MatchTo(to); x != tt.out {
			t.Errorf("Filter{To: %s}.MatchTo(%v) = %v, want %v", tt.to, to, x, tt.out)
		}
	}
}

func TestFilterMatchSubject(t *testing.T) {
	tests := []struct {
		subject string
		out     bool
	}{
		{"", false},
		{"Lorem", true},
		{"lorem", false}, // Case sensitivity is desired
		{"Lorum", false},
		{"Lorem ipsum dolor sit amet", true},
	}
	f := Filter{}
	subject := "Lorem ipsum dolor sit amet"
	for _, tt := range tests {
		f.Subject = tt.subject
		if x := f.MatchSubject(subject); x != tt.out {
			t.Errorf("Filter{Subject: %s}.MatchSubject(%s) = %v, want %v", tt.subject, subject, x, tt.out)
		}
	}
}

func TestFilterMatchOrigin(t *testing.T) {
	tests := []struct {
		origin string
		out    bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"127.0.0.2", false},
		{"10.0.0.1/24", false},
		{"127.0.0.1/24", true},
		{"127.0.0.1/32", true},
	}
	f := Filter{}
	addr := net.ParseIP("127.0.0.1")
	for _, tt := range tests {
		f.Origin = tt.origin
		if x := f.MatchOrigin(addr); x != tt.out {
			t.Errorf("Filter{Origin: %s}.MatchOrigin(%s) = %v, want %v", tt.origin, "127.0.0.1", x, tt.out)
		}
	}
}
