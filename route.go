package main

import (
	"sort"
)

type Route struct {
	Id        string
	Name      string
	To        string
	Hostname  string
	Port      int
	IsDefault bool
}

type RouteList []Route

// Implement sort.Interface
func (rl RouteList) Len() int {
	return len(rl)
}

func (rl RouteList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

// Ensure the DROP route is last.
func (rl RouteList) Less(i, j int) bool {
	if rl[i].Id == "DROP" {
		return false
	}
	if rl[j].Id == "DROP" {
		return true
	}
	return rl[i].Name < rl[j].Name
}

func SortedRoutes() RouteList {
	rl := make(RouteList, len(config.Routes))
	i := 0
	for _, route := range config.Routes {
		rl[i] = route
		i++
	}
	sort.Sort(rl)
	return rl
}
