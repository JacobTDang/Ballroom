package main

import "sort"

// FindItinerary reconstructs an itinerary from tickets, all starting from
// "JFK", using every ticket exactly once, returning the lexicographically
// smallest valid sequence of airports.
func FindItinerary(tickets [][]string) []string {
	graph := make(map[string][]string)
	for _, t := range tickets {
		graph[t[0]] = append(graph[t[0]], t[1])
	}
	for from := range graph {
		sort.Strings(graph[from])
	}

	var route []string
	var visit func(airport string)
	visit = func(airport string) {
		for len(graph[airport]) > 0 {
			next := graph[airport][0]
			graph[airport] = graph[airport][1:]
			visit(next)
		}
		route = append(route, airport)
	}
	visit("JFK")

	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}
	return route
}
