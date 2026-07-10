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

func toLevelOrder(root *TreeNode) []*int {
	if root == nil {
		return []*int{}
	}
	var out []*int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, nil)
			continue
		}
		v := node.Val
		out = append(out, &v)
		queue = append(queue, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == nil {
		out = out[:len(out)-1]
	}
	return out
}

func ip(v int) *int { return &v }

func TestSerializeDeserialize_RoundTrip(t *testing.T) {
	cases := [][]*int{
		{ip(1), ip(2), ip(3), nil, nil, ip(4), ip(5)},
		{},
		{ip(1)},
		{ip(-1), ip(-2), ip(-3)},
		{ip(5), ip(4), ip(7), ip(3), nil, ip(2), nil, ip(-1), nil, ip(9)},
	}

	for _, vals := range cases {
		original := buildTree(vals)
		roundTripped := Deserialize(Serialize(original))
		got := toLevelOrder(roundTripped)
		want := toLevelOrder(original)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("round trip of %v = %v, want %v", vals, got, want)
		}
	}
}
