package main

import "sort"

// Entry is one leaderboard row.
type Entry struct {
	Name  string
	Score int
}

// SortLeaderboard returns entries sorted by score descending; ties
// break by name ascending.
func SortLeaderboard(entries []Entry) []Entry {
	result := append([]Entry(nil), entries...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}
		return result[i].Name < result[j].Name
	})
	return result
}
