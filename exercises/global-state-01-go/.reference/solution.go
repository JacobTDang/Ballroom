package main

import "fmt"

// GenerateReport formats items (low-stock item names) into report
// lines, one per item.
func GenerateReport(items []string) []string {
	reportLines := make([]string, 0, len(items))
	for _, item := range items {
		reportLines = append(reportLines, fmt.Sprintf("LOW STOCK: %s", item))
	}
	return reportLines
}
