package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

type Filter struct {
	Id        string
	Order     int
	Name      string
	From      string
	To        string
	Subject   string
	Origin    string
	RouteId   string
	Summary   string // Convenience field for filter listing
	RouteName string // Convenience field for filter listing
}

func (f *Filter) Summarise() string {
	var attrs []string
	if f.From != "" {
		attrs = append(attrs, fmt.Sprintf("From: %s", f.From))
	}
	if f.To != "" {
		attrs = append(attrs, fmt.Sprintf("To: %s", f.To))
	}
	if f.Subject != "" {
		attrs = append(attrs, fmt.Sprintf("Subject: %s", f.Subject))
	}
	if f.Origin != "" {
		attrs = append(attrs, fmt.Sprintf("Origin: %s", f.Origin))
	}
	return strings.Join(attrs, ", ")
}

func (f *Filter) Match(from string, to []string, subject string, originIP net.IP) bool {
	fieldsSet := 0
	if f.From != "" {
		fieldsSet++
		if !f.MatchFrom(from) {
			return false
		}
	}
	if f.To != "" {
		fieldsSet++
		if !f.MatchTo(to) {
			return false
		}
	}
	if f.Subject != "" {
		fieldsSet++
		if !f.MatchSubject(subject) {
			return false
		}
	}
	if f.Origin != "" {
		fieldsSet++
		if !f.MatchOrigin(originIP) {
			return false
		}
	}
	// At this point all the fields that are set have been matched on.
	// Return false if none of the relevant fields are set, otherwise return true.
	return fieldsSet > 0
}

func (f *Filter) MatchFrom(from string) bool {
	if f.From == "" {
		return false
	}
	return strings.Contains(from, f.From)
}

func (f *Filter) MatchTo(to []string) bool {
	if f.To == "" {
		return false
	}
	// Test against all recipients
	for _, address := range to {
		if strings.Contains(address, f.To) {
			return true
		}
	}
	return false
}

func (f *Filter) MatchSubject(subject string) bool {
	if f.Subject == "" {
		return false
	}
	return strings.Contains(subject, f.Subject)
}

func (f *Filter) MatchOrigin(originIP net.IP) bool {
	// Is filter.Origin in CIDR notation e.g. "192.168.100.1/24" or "2001:DB8::/48"?
	_, filterNet, err := net.ParseCIDR(f.Origin)
	if err == nil {
		return filterNet.Contains(originIP)
	}

	// Is filter.Origin a valid IPv4 or IPv6 address?
	filterIP := net.ParseIP(f.Origin)
	if filterIP != nil {
		return filterIP.Equal(originIP)
	}

	return false
}

type FilterList []Filter

// Implement sort.Iterface
func (fl FilterList) Len() int {
	return len(fl)
}

func (fl FilterList) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}

func (fl FilterList) Less(i, j int) bool {
	return fl[i].Order < fl[j].Order
}

func SortedFilters() FilterList {
	fl := make(FilterList, len(config.Filters))
	i := 0
	for _, filter := range config.Filters {
		fl[i] = filter
		i++
	}
	sort.Sort(fl)
	return fl
}
