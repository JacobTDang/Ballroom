package main

import "testing"

func TestSettlesBill(t *testing.T) {
	cases := []struct {
		name    string
		amounts []float64
		bill    float64
		want    bool
	}{
		{"three dimes settle thirty cents", []float64{0.1, 0.1, 0.1}, 0.3, true},
		{"exact integer amounts settle", []float64{10.0, 20.0, 70.0}, 100.0, true},
		{"different rounding pattern settles", []float64{0.7, 0.1}, 0.8, true},
		{"another rounding pattern settles", []float64{1.1, 2.2}, 3.3, true},
		{"short by a cent does not settle", []float64{10.00, 10.00}, 20.01, false},
		{"clearly short does not settle", []float64{5.00, 5.00}, 11.00, false},
		{"empty amounts settle zero bill", []float64{}, 0.0, true},
	}
	for _, c := range cases {
		if got := SettlesBill(c.amounts, c.bill); got != c.want {
			t.Errorf("%s: SettlesBill(%v, %v) = %v, want %v", c.name, c.amounts, c.bill, got, c.want)
		}
	}
}
