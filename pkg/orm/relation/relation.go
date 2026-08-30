// Package relation defines the concrete, strongly-typed relationship types
// used to declare relationship structures and compile nested sub-queries.
package relation

import (
	"context"
	"errors"
	"fmt"

	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

type nestedQueryState struct {
	wheres []contract.Where
}

func (n nestedQueryState) GetSchemaName() string { return "" }
func (n nestedQueryState) GetTableName() string  { return "" }
func (n nestedQueryState) GetDefaultColumns() []string {
	return nil
}
func (n nestedQueryState) GetSelectedColumns() []string {
	return nil
}
func (n nestedQueryState) IsDistinct() bool { return false }
func (n nestedQueryState) GetAggregate() *contract.AggregateState {
	return nil
}
func (n nestedQueryState) GetWheres() []contract.Where { return n.wheres }
func (n nestedQueryState) GetOrders() []contract.Order { return nil }
func (n nestedQueryState) GetGroups() []string         { return nil }
func (n nestedQueryState) GetHavings() []contract.Where {
	return nil
}
func (n nestedQueryState) GetLimit() (uint64, bool)  { return 0, false }
func (n nestedQueryState) GetOffset() (uint64, bool) { return 0, false }
func (n nestedQueryState) GetColumnCastAndTypedSlice(
	_ string,
	slice []any,
) (string, any, bool) {
	return "", slice, false
}

// Relation is the base concrete generic type representing any database
// relationship.
type Relation[
	ParentModel any,
	ChildModel any,
	ChildModelMutator contract.Mutator,
] struct {
	QueryBuilder      *orm.QueryBuilder[ChildModel, ChildModelMutator]
	ForeignKeyColumns []contract.Column[ChildModel]
	ChildQueryFactory func(
		db contract.DB,
	) *orm.QueryBuilder[ChildModel, ChildModelMutator]
	LocalKeyExtractor   func(parent *ParentModel) any
	ForeignKeyExtractor func(child *ChildModel) any
	Hydrator            func(parent *ParentModel, children []ChildModel)
	Customizers []func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator]
}

// Customise appends a query customizer to the relationship.
func (r Relation[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Customise(
	customizer func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator],
) Relation[ParentModel, ChildModel, ChildModelMutator] {
	r.Customizers = append(
		append(
			[]func(
				qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
			) *orm.QueryBuilder[ChildModel, ChildModelMutator](nil),
			r.Customizers...,
		),
		customizer,
	)
	return r
}

// GetRelationQuery returns the underlying child sub-query state provider.
func (r Relation[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) GetRelationQuery() contract.QueryStateProvider {
	var qb *orm.QueryBuilder[ChildModel, ChildModelMutator]
	if r.ChildQueryFactory != nil {
		qb = r.ChildQueryFactory(nil)
	} else if r.QueryBuilder != nil {
		qb = r.QueryBuilder.Clone()
	} else {
		return nil
	}

	for _, customizer := range r.Customizers {
		qb = customizer(qb)
	}
	return qb
}

// EagerLoad loads relation records and assigns them to parents.
func (r Relation[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) EagerLoad(
	ctxEagerLoad context.Context,
	db contract.DB,
	parents []ParentModel,
) error {
	if len(parents) == 0 {
		return nil
	}

	localKeys := make([]any, len(parents))
	for i := range parents {
		localKeys[i] = r.LocalKeyExtractor(&parents[i])
	}

	var qb *orm.QueryBuilder[ChildModel, ChildModelMutator]
	if r.ChildQueryFactory != nil {
		qb = r.ChildQueryFactory(db)
	} else if r.QueryBuilder != nil {
		qb = r.QueryBuilder.Clone()
		qb.SetDB(db)
	} else {
		return errors.New("relation has no query builder or factory")
	}

	for _, customizer := range r.Customizers {
		qb = customizer(qb)
	}

	switch len(r.ForeignKeyColumns) {
	case 0:
		return errors.New("relation has no foreign key columns defined")
	case 1:
		qb.Where(where.InWhere{
			Column:  r.ForeignKeyColumns[0].ColumnName(),
			Values:  localKeys,
			Boolean: contract.BoolAnd,
		})
	default:
		var compositeWheres []contract.Where
		for _, key := range localKeys {
			parts, ok := key.([]any)
			if !ok {
				continue
			}
			var andWheres []contract.Where
			for idx, col := range r.ForeignKeyColumns {
				andWheres = append(andWheres, where.BasicWhere{
					Column:   col.ColumnName(),
					Operator: contract.OpEqual,
					Value:    parts[idx],
					Boolean:  contract.BoolAnd,
				})
			}
			compositeWheres = append(compositeWheres, where.NestedWhere{
				Query:   nestedQueryState{wheres: andWheres},
				Boolean: contract.BoolOr,
			})
		}
		qb.Where(where.NestedWhere{
			Query:   nestedQueryState{wheres: compositeWheres},
			Boolean: contract.BoolAnd,
		})
	}

	children, err := qb.Get(ctxEagerLoad)
	if err != nil {
		return fmt.Errorf("failed executing relation query -> %w", err)
	}

	mapParentKeyToChildren := make(map[any][]ChildModel)
	for i := range children {
		key := toComparableKey(r.ForeignKeyExtractor(&children[i]))
		mapParentKeyToChildren[key] = append(
			mapParentKeyToChildren[key],
			children[i],
		)
	}

	for i := range parents {
		key := toComparableKey(r.LocalKeyExtractor(&parents[i]))
		r.Hydrator(&parents[i], mapParentKeyToChildren[key])
	}

	return nil
}

func toComparableKey(val any) any {
	slice, ok := val.([]any)
	if ok {
		derefed := make([]any, len(slice))
		for i, item := range slice {
			derefed[i] = derefValue(item)
		}
		return fmt.Sprintf("%v", derefed)
	}
	return derefValue(val)
}

func derefValue(val any) any {
	if val == nil {
		return nil
	}
	switch p := val.(type) {
	case *int:
		if p == nil {
			return nil
		}
		return *p
	case *int64:
		if p == nil {
			return nil
		}
		return *p
	case *int32:
		if p == nil {
			return nil
		}
		return *p
	case *int16:
		if p == nil {
			return nil
		}
		return *p
	case *int8:
		if p == nil {
			return nil
		}
		return *p
	case *string:
		if p == nil {
			return nil
		}
		return *p
	case *bool:
		if p == nil {
			return nil
		}
		return *p
	case *float64:
		if p == nil {
			return nil
		}
		return *p
	case *float32:
		if p == nil {
			return nil
		}
		return *p
	}
	return val
}
