package relation

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
)

// HasMany represents a 1:N relationship where the parent model has many
// child models.
type HasMany[
	ParentModel any,
	ChildModel any,
	ChildModelMutator contract.Mutator,
] struct {
	Relation[ParentModel, ChildModel, ChildModelMutator]
}

// Where allows filtering the relation's nested query builder inline while
// preserving exact types and avoiding query filter leakage via state
// cloning.
func (r HasMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Where(
	w contract.Where,
) HasMany[ParentModel, ChildModel, ChildModelMutator] {
	r.Relation = r.Relation.Customise(func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator] {
		return qb.Where(w)
	})
	return r
}

// Constrain allows customizing the relation's nested query builder inline.
func (r HasMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Constrain(
	callback func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator],
) HasMany[ParentModel, ChildModel, ChildModelMutator] {
	r.Relation = r.Relation.Customise(callback)
	return r
}
