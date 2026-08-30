package aip160

import (
	"go.chromium.org/luci/common/data/aip160"
	"slices"
)

// WalkFunc is a callback function applied to each node during tree traversal.
type WalkFunc func(node any) error

// IterativeWalk traverses a tree-like structure in post-order using an
// iterative approach and applies a callback on each node.
func IterativeWalk(root *aip160.Filter, fn WalkFunc) error {
	if root.Expression == nil {
		return nil
	}

	workStack := []any{root.Expression}
	postOrderStack := make([]any, 0, 32)

	for len(workStack) > 0 {
		curr := workStack[len(workStack)-1]
		workStack = workStack[:len(workStack)-1]

		postOrderStack = append(postOrderStack, curr)

		switch n := curr.(type) {
		case *aip160.Expression:
			for _, v := range slices.Backward(n.Sequences) {
				workStack = append(workStack, v)
			}
		case *aip160.Sequence:
			for _, v := range slices.Backward(n.Factors) {
				workStack = append(workStack, v)
			}
		case *aip160.Factor:
			for _, v := range slices.Backward(n.Terms) {
				workStack = append(workStack, v)
			}
		case *aip160.Term:
			if n.Simple != nil {
				workStack = append(workStack, n.Simple)
			}
		case *aip160.Simple:
			if n.Restriction != nil {
				workStack = append(workStack, n.Restriction)
			} else if n.Composite != nil {
				workStack = append(workStack, n.Composite)
			}
		}
	}

	for _, p := range slices.Backward(postOrderStack) {
		err := fn(p)
		if err != nil {
			return err
		}
	}

	return nil
}
