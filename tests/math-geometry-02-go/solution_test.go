package main

import (
	"reflect"
	"testing"
)

func TestSpiralOrder_3x3(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	want := []int{1, 2, 3, 6, 9, 8, 7, 4, 5}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_3x4(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	}
	want := []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_SingleRow(t *testing.T) {
	matrix := [][]int{{1, 2, 3, 4}}
	want := []int{1, 2, 3, 4}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_SingleColumn(t *testing.T) {
	matrix := [][]int{{1}, {2}, {3}}
	want := []int{1, 2, 3}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_SingleElement(t *testing.T) {
	matrix := [][]int{{7}}
	want := []int{7}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_4x3(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	}
	want := []int{1, 2, 3, 6, 9, 12, 11, 10, 7, 4, 5, 8}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_2x2(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}
	want := []int{1, 2, 4, 3}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_NegativeValues(t *testing.T) {
	matrix := [][]int{
		{-1, -2},
		{-3, -4},
	}
	want := []int{-1, -2, -4, -3}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}

func TestSpiralOrder_4x4(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	want := []int{1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10}
	if got := SpiralOrder(matrix); !reflect.DeepEqual(got, want) {
		t.Errorf("SpiralOrder(%v) = %v, want %v", matrix, got, want)
	}
}
