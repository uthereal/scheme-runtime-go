package relation

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
)

// BelongsTo represents an N:1 relationship where the parent model belongs
// to a child model.
type BelongsTo[
	ParentModel any,
	ChildModel any,
	ChildModelMutator contract.Mutator,
] struct {
	Relation[ParentModel, ChildModel, ChildModelMutator]
}

// Where allows filtering the relation's nested query builder inline while
// preserving exact types and avoiding query filter leakage via state
// cloning.
func (r BelongsTo[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Where(
	w contract.Where,
) BelongsTo[ParentModel, ChildModel, ChildModelMutator] {
	r.Relation = r.Relation.Customise(func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator] {
		return qb.Where(w)
	})
	return r
}

// Constrain allows customizing the relation's nested query builder inline.
func (r BelongsTo[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Constrain(
	callback func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator],
) BelongsTo[ParentModel, ChildModel, ChildModelMutator] {
	r.Relation = r.Relation.Customise(callback)
	return r
}
