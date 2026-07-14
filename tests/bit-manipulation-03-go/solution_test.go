package main

import (
	"reflect"
	"testing"
)

func TestCountBits_Classic(t *testing.T) {
	want := []int{0, 1, 1, 2, 1}
	if got := CountBits(4); !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(4) = %v, want %v", got, want)
	}
}

func TestCountBits_Small(t *testing.T) {
	want := []int{0, 1, 1}
	if got := CountBits(2); !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(2) = %v, want %v", got, want)
	}
}

func TestCountBits_Zero(t *testing.T) {
	want := []int{0}
	if got := CountBits(0); !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(0) = %v, want %v", got, want)
	}
}

func TestCountBits_One(t *testing.T) {
	want := []int{0, 1}
	if got := CountBits(1); !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(1) = %v, want %v", got, want)
	}
}

func TestCountBits_Larger(t *testing.T) {
	want := []int{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}
	if got := CountBits(15); !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(15) = %v, want %v", got, want)
	}
}
