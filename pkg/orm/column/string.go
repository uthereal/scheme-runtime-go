package column

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// StringColumn represents a strongly-typed column schema for string-based
// fields.
type StringColumn[Model any, Type ~string] struct {
	Column[Model, Type]
}

// Contains returns a condition checking if this column contains the
// specified substring (e.g., LIKE %val%).
func (c StringColumn[Model, Type]) Contains(val string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLike,
		Value:    "%" + val + "%",
		Boolean:  contract.BoolAnd,
	}
}

// ContainsFold returns a case-insensitive condition checking if this
// column contains the substring.
func (c StringColumn[Model, Type]) ContainsFold(val string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpILike,
		Value:    "%" + val + "%",
		Boolean:  contract.BoolAnd,
	}
}

// HasPrefix returns a condition checking if this column starts with the
// specified prefix (e.g., LIKE val%).
func (c StringColumn[Model, Type]) HasPrefix(val string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLike,
		Value:    val + "%",
		Boolean:  contract.BoolAnd,
	}
}

// HasSuffix returns a condition checking if this column ends with the
// specified suffix (e.g., LIKE %val).
func (c StringColumn[Model, Type]) HasSuffix(val string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLike,
		Value:    "%" + val,
		Boolean:  contract.BoolAnd,
	}
}

// Like returns a SQL pattern-matching condition (e.g., column LIKE
// pattern).
func (c StringColumn[Model, Type]) Like(pattern string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpLike,
		Value:    pattern,
		Boolean:  contract.BoolAnd,
	}
}

// ILike returns a case-insensitive SQL pattern-matching condition (e.g.,
// column ILIKE pattern).
func (c StringColumn[Model, Type]) ILike(pattern string) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpILike,
		Value:    pattern,
		Boolean:  contract.BoolAnd,
	}
}

// TextSearchMatch returns a full-text search matching condition against the
// query string (e.g., column @@ query).
func (c StringColumn[Model, Type]) TextSearchMatch(
	query string,
) contract.Where {
	return where.BasicWhere{
		Column:   c.Name,
		Operator: contract.OpTextSearchMatch,
		Value:    query,
		Boolean:  contract.BoolAnd,
	}
}
