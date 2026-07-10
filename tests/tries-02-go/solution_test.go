package main

import "testing"

func TestWordDictionary(t *testing.T) {
	d := NewWordDictionary()
	d.AddWord("bad")
	d.AddWord("dad")
	d.AddWord("mad")

	cases := []struct {
		query string
		want  bool
	}{
		{"pad", false},
		{"bad", true},
		{".ad", true},
		{"b..", true},
		{"...", true},
		{"....", false},
		{"..d", true},
		{"dab", false},
	}

	for _, c := range cases {
		if got := d.Search(c.query); got != c.want {
			t.Errorf("Search(%q) = %v, want %v", c.query, got, c.want)
		}
	}
}

func TestWordDictionary_EmptyDictionaryNeverMatches(t *testing.T) {
	d := NewWordDictionary()
	if d.Search("a") {
		t.Error("Search on empty dictionary should be false")
	}
	if d.Search(".") {
		t.Error("Search('.') on empty dictionary should be false")
	}
}
