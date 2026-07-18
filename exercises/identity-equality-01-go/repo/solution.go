package main

// Record is a deduplicable value: two records are duplicates if they
// have the same Key and Value.
type Record struct {
	Key   string
	Value int
}

// Dedupe removes duplicate records -- two records with the same Key
// and Value are duplicates. Currently keeps both -- find and fix the
// bug.
func Dedupe(records []*Record) []*Record {
	seen := map[*Record]bool{}
	var result []*Record
	for _, r := range records {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return result
}
