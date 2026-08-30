package contract

import "golang.org/x/exp/constraints"

// Number represents any integer or floating-point numeric type.
type Number interface {
	constraints.Integer | constraints.Float
}
