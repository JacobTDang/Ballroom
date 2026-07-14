package main

import "testing"

func TestCanAttendMeetings(t *testing.T) {
	cases := []struct {
		name      string
		intervals [][]int
		want      bool
	}{
		{
			name:      "overlap in the middle",
			intervals: [][]int{{0, 30}, {5, 10}, {15, 20}},
			want:      false,
		},
		{
			name:      "no overlap",
			intervals: [][]int{{7, 10}, {2, 4}},
			want:      true,
		},
		{
			name:      "touching endpoints is fine",
			intervals: [][]int{{5, 10}, {10, 15}},
			want:      true,
		},
		{
			name:      "empty schedule",
			intervals: [][]int{},
			want:      true,
		},
		{
			name:      "single meeting",
			intervals: [][]int{{3, 8}},
			want:      true,
		},
		{
			name:      "unsorted overlap at the end",
			intervals: [][]int{{13, 15}, {1, 5}, {6, 8}, {14, 20}},
			want:      false,
		},
		{
			name:      "boundary values, overlap by one",
			intervals: [][]int{{0, 1000000}, {999999, 1000000}},
			want:      false,
		},
		{
			name:      "larger schedule, all sequential",
			intervals: [][]int{{0, 10}, {10, 20}, {20, 30}, {30, 40}, {40, 50}, {50, 60}, {60, 70}},
			want:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CanAttendMeetings(c.intervals); got != c.want {
				t.Errorf("CanAttendMeetings(%v) = %v, want %v", c.intervals, got, c.want)
			}
		})
	}
}
