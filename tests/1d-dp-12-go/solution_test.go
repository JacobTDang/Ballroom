package main

import "testing"

func TestCanPartition_Classic(t *testing.T) {
	nums := []int{1, 5, 11, 5}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}

func TestCanPartition_OddSum(t *testing.T) {
	nums := []int{1, 2, 3, 5}
	if got := CanPartition(nums); got != false {
		t.Errorf("CanPartition(%v) = %v, want false", nums, got)
	}
}

func TestCanPartition_EvenSplit(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}

func TestCanPartition_TwoEqual(t *testing.T) {
	nums := []int{2, 2}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}

func TestCanPartition_SingleElement(t *testing.T) {
	nums := []int{4}
	if got := CanPartition(nums); got != false {
		t.Errorf("CanPartition(%v) = %v, want false", nums, got)
	}
}

func TestCanPartition_AllSame(t *testing.T) {
	nums := []int{3, 3, 3, 3}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}

func TestCanPartition_EvenSumUnreachable(t *testing.T) {
	nums := []int{2, 2, 3, 5}
	if got := CanPartition(nums); got != false {
		t.Errorf("CanPartition(%v) = %v, want false", nums, got)
	}
}

func TestCanPartition_BoundaryValues(t *testing.T) {
	nums := []int{100, 100, 100, 100}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}

func TestCanPartition_LargerMultiCombination(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	if got := CanPartition(nums); got != true {
		t.Errorf("CanPartition(%v) = %v, want true", nums, got)
	}
}
