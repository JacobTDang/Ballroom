package main

import (
	"reflect"
	"testing"
)

func TestRotateImage_3x3(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	want := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_2x2(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}
	want := [][]int{
		{3, 1},
		{4, 2},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_1x1(t *testing.T) {
	matrix := [][]int{{5}}
	want := [][]int{{5}}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_4x4(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	want := [][]int{
		{13, 9, 5, 1},
		{14, 10, 6, 2},
		{15, 11, 7, 3},
		{16, 12, 8, 4},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_NegativeValues(t *testing.T) {
	matrix := [][]int{
		{-1, -2},
		{-3, -4},
	}
	want := [][]int{
		{-3, -1},
		{-4, -2},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_AllSameValues(t *testing.T) {
	matrix := [][]int{
		{7, 7, 7},
		{7, 7, 7},
		{7, 7, 7},
	}
	want := [][]int{
		{7, 7, 7},
		{7, 7, 7},
		{7, 7, 7},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_WithZero(t *testing.T) {
	matrix := [][]int{
		{0, 1},
		{2, 3},
	}
	want := [][]int{
		{2, 0},
		{3, 1},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}

func TestRotateImage_5x5(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		{11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20},
		{21, 22, 23, 24, 25},
	}
	want := [][]int{
		{21, 16, 11, 6, 1},
		{22, 17, 12, 7, 2},
		{23, 18, 13, 8, 3},
		{24, 19, 14, 9, 4},
		{25, 20, 15, 10, 5},
	}
	RotateImage(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("RotateImage() = %v, want %v", matrix, want)
	}
}
