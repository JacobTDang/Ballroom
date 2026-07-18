package main

import "fmt"

// reportLines accumulates formatted report lines. It's meant to be
// scratch space for one call to GenerateReport, not storage that
// outlives it.
var reportLines []string

// GenerateReport formats items (low-stock item names) into report
// lines, one per item. Currently a call's report can contain lines
// left over from an earlier call -- find and fix the bug.
func GenerateReport(items []string) []string {
	for _, item := range items {
		reportLines = append(reportLines, fmt.Sprintf("LOW STOCK: %s", item))
	}
	return reportLines
}
