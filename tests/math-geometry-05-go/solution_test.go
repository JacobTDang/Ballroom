package main

import (
	"reflect"
	"testing"
)

func TestPlusOne_Simple(t *testing.T) {
	got := PlusOne([]int{1, 2, 3})
	want := []int{1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([1,2,3]) = %v, want %v", got, want)
	}
}

func TestPlusOne_AllNines(t *testing.T) {
	got := PlusOne([]int{9, 9, 9})
	want := []int{1, 0, 0, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([9,9,9]) = %v, want %v", got, want)
	}
}

func TestPlusOne_SingleZero(t *testing.T) {
	got := PlusOne([]int{0})
	want := []int{1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([0]) = %v, want %v", got, want)
	}
}

func TestPlusOne_TrailingNine(t *testing.T) {
	got := PlusOne([]int{1, 2, 9})
	want := []int{1, 3, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([1,2,9]) = %v, want %v", got, want)
	}
}

func TestPlusOne_SingleNine(t *testing.T) {
	got := PlusOne([]int{9})
	want := []int{1, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([9]) = %v, want %v", got, want)
	}
}

func TestPlusOne_PartialTrailingNines(t *testing.T) {
	got := PlusOne([]int{1, 9, 9})
	want := []int{2, 0, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([1,9,9]) = %v, want %v", got, want)
	}
}

func TestPlusOne_MixedNoCarryPastStop(t *testing.T) {
	got := PlusOne([]int{9, 8, 9, 9})
	want := []int{9, 9, 0, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([9,8,9,9]) = %v, want %v", got, want)
	}
}

func TestPlusOne_SingleDigitNotNine(t *testing.T) {
	got := PlusOne([]int{5})
	want := []int{6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([5]) = %v, want %v", got, want)
	}
}

func TestPlusOne_LargerAllNines(t *testing.T) {
	got := PlusOne([]int{9, 9, 9, 9, 9})
	want := []int{1, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlusOne([9,9,9,9,9]) = %v, want %v", got, want)
	}
}
