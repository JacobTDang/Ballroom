package main

import (
	"reflect"
	"testing"
)

func TestPartitionLabels_Classic(t *testing.T) {
	s := "ababcbacadefegdehijhklij"
	want := []int{9, 7, 8}
	if got := PartitionLabels(s); !reflect.DeepEqual(got, want) {
		t.Errorf("PartitionLabels(%q) = %v, want %v", s, got, want)
	}
}

func TestPartitionLabels_AllUnique(t *testing.T) {
	s := "abcde"
	want := []int{1, 1, 1, 1, 1}
	if got := PartitionLabels(s); !reflect.DeepEqual(got, want) {
		t.Errorf("PartitionLabels(%q) = %v, want %v", s, got, want)
	}
}

func TestPartitionLabels_AllSame(t *testing.T) {
	s := "aaaa"
	want := []int{4}
	if got := PartitionLabels(s); !reflect.DeepEqual(got, want) {
		t.Errorf("PartitionLabels(%q) = %v, want %v", s, got, want)
	}
}

func TestPartitionLabels_SingleChar(t *testing.T) {
	s := "a"
	want := []int{1}
	if got := PartitionLabels(s); !reflect.DeepEqual(got, want) {
		t.Errorf("PartitionLabels(%q) = %v, want %v", s, got, want)
	}
}

func TestPartitionLabels_MultipleEqualPartitions(t *testing.T) {
	s := "aabbcc"
	want := []int{2, 2, 2}
	if got := PartitionLabels(s); !reflect.DeepEqual(got, want) {
		t.Errorf("PartitionLabels(%q) = %v, want %v", s, got, want)
	}
}
