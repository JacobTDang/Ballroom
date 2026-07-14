package main

import "testing"

func TestMinMeetingRooms(t *testing.T) {
	cases := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{
			name:      "classic two rooms",
			intervals: [][]int{{0, 30}, {5, 10}, {15, 20}},
			want:      2,
		},
		{
			name:      "no overlap, one room",
			intervals: [][]int{{7, 10}, {2, 4}},
			want:      1,
		},
		{
			name:      "touching endpoints share a room",
			intervals: [][]int{{5, 10}, {10, 15}},
			want:      1,
		},
		{
			name:      "empty schedule",
			intervals: [][]int{},
			want:      0,
		},
		{
			name:      "three identical meetings need three rooms",
			intervals: [][]int{{1, 2}, {1, 2}, {1, 2}},
			want:      3,
		},
		{
			name:      "single meeting",
			intervals: [][]int{{1, 5}},
			want:      1,
		},
		{
			name:      "five fully overlapping meetings need five rooms",
			intervals: [][]int{{1, 100}, {1, 100}, {1, 100}, {1, 100}, {1, 100}},
			want:      5,
		},
		{
			name:      "staggered overlaps, larger input",
			intervals: [][]int{{1, 10}, {2, 7}, {3, 19}, {8, 12}, {10, 20}, {11, 30}},
			want:      4,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := MinMeetingRooms(c.intervals); got != c.want {
				t.Errorf("MinMeetingRooms(%v) = %d, want %d", c.intervals, got, c.want)
			}
		})
	}
}
