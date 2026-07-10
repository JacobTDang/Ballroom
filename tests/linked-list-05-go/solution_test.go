package main

import "testing"

// buildRandomList builds a linked list from vals, wiring each node's
// Random pointer to the node at randomIdx[i] (or nil for -1).
func buildRandomList(vals []int, randomIdx []int) []*Node {
	nodes := make([]*Node, len(vals))
	for i, v := range vals {
		nodes[i] = &Node{Val: v}
	}
	for i := range nodes {
		if i+1 < len(nodes) {
			nodes[i].Next = nodes[i+1]
		}
		if randomIdx[i] >= 0 {
			nodes[i].Random = nodes[randomIdx[i]]
		}
	}
	return nodes
}

func checkCopyRandomList(t *testing.T, vals []int, randomIdx []int) {
	t.Helper()
	orig := buildRandomList(vals, randomIdx)
	origSet := make(map[*Node]bool, len(orig))
	for _, n := range orig {
		origSet[n] = true
	}

	var head *Node
	if len(orig) > 0 {
		head = orig[0]
	}
	copyHead := CopyRandomList(head)

	var copies []*Node
	for cur := copyHead; cur != nil; cur = cur.Next {
		if origSet[cur] {
			t.Fatal("copied node shares identity with an original node — not a deep copy")
		}
		copies = append(copies, cur)
	}
	if len(copies) != len(vals) {
		t.Fatalf("copied list length = %d, want %d", len(copies), len(vals))
	}

	copyIndex := make(map[*Node]int, len(copies))
	for i, n := range copies {
		copyIndex[n] = i
	}

	for i, n := range copies {
		if n.Val != vals[i] {
			t.Errorf("copies[%d].Val = %d, want %d", i, n.Val, vals[i])
		}
		if randomIdx[i] == -1 {
			if n.Random != nil {
				t.Errorf("copies[%d].Random = non-nil, want nil", i)
			}
			continue
		}
		if n.Random == nil {
			t.Errorf("copies[%d].Random = nil, want index %d", i, randomIdx[i])
			continue
		}
		gotIdx, ok := copyIndex[n.Random]
		if !ok || gotIdx != randomIdx[i] {
			t.Errorf("copies[%d].Random points to index %v (ok=%v), want %d", i, gotIdx, ok, randomIdx[i])
		}
	}
}

func TestCopyRandomList(t *testing.T) {
	checkCopyRandomList(t, []int{7, 13, 11, 10, 1}, []int{-1, 0, 4, 2, 0})
	checkCopyRandomList(t, []int{1, 2}, []int{1, 1})
	checkCopyRandomList(t, []int{3, 3, 3}, []int{-1, -1, -1})
	checkCopyRandomList(t, []int{}, []int{})
}
