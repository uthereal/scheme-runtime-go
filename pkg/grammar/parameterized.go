package grammar

import "strconv"

// bindingsTracker tracks all generated positional parameter values for
// prepared statements.
type bindingsTracker struct {
	values []any
}

// Bind returns a placeholder string and a brand-new, independent
// bindingsTracker instance carrying the appended value, strictly
// avoiding any in-place mutation of the original tracker.
func (b bindingsTracker) Bind(val any) (string, bindingsTracker) {
	newValues := make([]any, len(b.values), len(b.values)+1)
	copy(newValues, b.values)
	newValues = append(newValues, val)

	placeholder := "$" + strconv.Itoa(len(newValues))
	return placeholder, bindingsTracker{values: newValues}
}
