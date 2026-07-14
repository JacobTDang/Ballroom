package main

import (
	"reflect"
	"testing"
)

func buildTree(vals []*int) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: *vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != nil {
				node.Left = &TreeNode{Val: *vals[i]}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				node.Right = &TreeNode{Val: *vals[i]}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

func ip(v int) *int { return &v }

func TestLevelOrder(t *testing.T) {
	cases := []struct {
		in   []*int
		want [][]int
	}{
		{
			[]*int{ip(3), ip(9), ip(20), nil, nil, ip(15), ip(7)},
			[][]int{{3}, {9, 20}, {15, 7}},
		},
		{[]*int{ip(1)}, [][]int{{1}}},
		{[]*int{}, nil},
		{[]*int{ip(1), ip(2), ip(3), ip(4), ip(5), ip(6), ip(7)}, [][]int{{1}, {2, 3}, {4, 5, 6, 7}}},
		{[]*int{ip(1), ip(2)}, [][]int{{1}, {2}}},
	}

	for _, c := range cases {
		got := LevelOrder(buildTree(c.in))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("LevelOrder(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
