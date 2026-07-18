package main

import (
	"reflect"
	"testing"
)

func names(entries []Entry) []string {
	var out []string
	for _, e := range entries {
		out = append(out, e.Name)
	}
	return out
}

func TestSortLeaderboard(t *testing.T) {
	cases := []struct {
		name    string
		entries []Entry
		want    []string
	}{
		{
			"no ties sorts by score only",
			[]Entry{{"bob", 60}, {"dan", 100}, {"cara", 75}, {"amy", 90}},
			[]string{"dan", "amy", "cara", "bob"},
		},
		{
			"tie group breaks ascending by name",
			[]Entry{{"erin", 90}, {"cara", 75}, {"amy", 90}, {"bob", 90}},
			[]string{"amy", "bob", "erin", "cara"},
		},
		{
			"multiple tie groups",
			[]Entry{{"zoe", 50}, {"amy", 80}, {"erin", 65}, {"dan", 50}, {"cara", 80}, {"bob", 80}},
			[]string{"amy", "bob", "cara", "erin", "dan", "zoe"},
		},
		{
			"fully tied list sorts by name",
			[]Entry{{"zed", 10}, {"amy", 10}, {"mno", 10}},
			[]string{"amy", "mno", "zed"},
		},
		{
			"negative scores with tie",
			[]Entry{{"cara", 10}, {"bob", -5}, {"amy", -5}},
			[]string{"cara", "amy", "bob"},
		},
	}
	for _, c := range cases {
		got := names(SortLeaderboard(c.entries))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
