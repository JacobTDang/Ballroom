package main

import (
	"reflect"
	"testing"
)

func TestSetZeroes_Classic(t *testing.T) {
	matrix := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	want := [][]int{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_TwoZeroes(t *testing.T) {
	matrix := [][]int{
		{0, 1, 2, 0},
		{3, 4, 5, 2},
		{1, 3, 1, 5},
	}
	want := [][]int{
		{0, 0, 0, 0},
		{0, 4, 5, 0},
		{0, 3, 1, 0},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_SingleZero(t *testing.T) {
	matrix := [][]int{
		{1, 0},
		{1, 1},
	}
	want := [][]int{
		{0, 0},
		{1, 0},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_NoZero(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}
	want := [][]int{
		{1, 2},
		{3, 4},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_AllZeros(t *testing.T) {
	matrix := [][]int{
		{0, 0},
		{0, 0},
	}
	want := [][]int{
		{0, 0},
		{0, 0},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_CornerZero(t *testing.T) {
	matrix := [][]int{
		{0, 1},
		{1, 1},
	}
	want := [][]int{
		{0, 0},
		{0, 1},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_SingleRow(t *testing.T) {
	matrix := [][]int{{1, 0, 3}}
	want := [][]int{{0, 0, 0}}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_SingleColumn(t *testing.T) {
	matrix := [][]int{{1}, {0}, {3}}
	want := [][]int{{0}, {0}, {0}}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}

func TestSetZeroes_NegativeValues(t *testing.T) {
	matrix := [][]int{
		{-1, 0},
		{-2, -3},
	}
	want := [][]int{
		{0, 0},
		{-2, 0},
	}
	SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, want) {
		t.Errorf("SetZeroes() = %v, want %v", matrix, want)
	}
}
