package main

import (
	"reflect"
	"testing"
)

// One sequence of three calls exercises all three traps: a lone call
// passes even with the bug (nothing to leak from yet), a second call
// with different input is polluted by the first, and a third call
// repeating the first call's input shows that pollution as an
// outright duplicate.
func TestGenerateReport_CallsDoNotLeakIntoEachOther(t *testing.T) {
	cases := []struct {
		items []string
		want  []string
	}{
		{[]string{"apples"}, []string{"LOW STOCK: apples"}},
		{[]string{"bananas"}, []string{"LOW STOCK: bananas"}},
		{[]string{"apples"}, []string{"LOW STOCK: apples"}},
	}
	for _, c := range cases {
		if got := GenerateReport(c.items); !reflect.DeepEqual(got, c.want) {
			t.Errorf("GenerateReport(%v) = %v, want %v", c.items, got, c.want)
		}
	}
}
