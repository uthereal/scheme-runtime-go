package relation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// BelongsToMany represents an N:N relationship linked via a pivot table.
type BelongsToMany[
	ParentModel any,
	ChildModel any,
	ChildModelMutator contract.Mutator,
] struct {
	PivotSchema             string
	PivotTable              string
	PivotForeignKeyToParent string
	PivotForeignKeyToChild  string
	ChildForeignKeyColumn   contract.Column[ChildModel]
	ChildQueryFactory       func(
		db contract.DB,
	) *orm.QueryBuilder[ChildModel, ChildModelMutator]
	LocalKeyExtractor func(parent *ParentModel) any
	ChildKeyExtractor func(child *ChildModel) any
	Hydrator          func(parent *ParentModel, children []ChildModel)
	Customizers       []func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator]
}

// Customise appends a query customizer to the relationship.
func (r BelongsToMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Customise(
	customizer func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator],
) BelongsToMany[ParentModel, ChildModel, ChildModelMutator] {
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

// Where allows filtering the relation's nested query builder inline while
// preserving exact types and avoiding query filter leakage via state
// cloning.
func (r BelongsToMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Where(
	w contract.Where,
) BelongsToMany[ParentModel, ChildModel, ChildModelMutator] {
	return r.Customise(func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator] {
		return qb.Where(w)
	})
}

// Constrain allows customizing the relation's nested query builder inline.
func (r BelongsToMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) Constrain(
	callback func(
		qb *orm.QueryBuilder[ChildModel, ChildModelMutator],
	) *orm.QueryBuilder[ChildModel, ChildModelMutator],
) BelongsToMany[ParentModel, ChildModel, ChildModelMutator] {
	return r.Customise(callback)
}

// GetRelationQuery returns the underlying child sub-query state provider.
func (r BelongsToMany[
	ParentModel,
	ChildModel,
	ChildModelMutator,
]) GetRelationQuery() contract.QueryStateProvider {
	var qb *orm.QueryBuilder[ChildModel, ChildModelMutator]
	if r.ChildQueryFactory != nil {
		rawQb := r.ChildQueryFactory(nil)
		if rawQb != nil {
			qb = rawQb.Clone()
		}
	} else {
		return nil
	}

	for _, customizer := range r.Customizers {
		qb = customizer(qb)
	}
	return qb
}

// EagerLoad loads relation records via pivot table and assigns them to parents.
func (r BelongsToMany[
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

	if r.LocalKeyExtractor == nil {
		return errors.New("relation has no local key extractor")
	}
	if r.ChildQueryFactory == nil {
		return errors.New("relation has no child query factory")
	}
	if r.ChildForeignKeyColumn == nil {
		return errors.New("relation has no child foreign key column")
	}
	if r.ChildKeyExtractor == nil {
		return errors.New("relation has no child key extractor")
	}
	if r.Hydrator == nil {
		return errors.New("relation has no hydrator")
	}

	localKeys := make([]any, 0, len(parents))
	mapLocalKeyToIsSeen := make(map[any]bool)
	for i := range parents {
		rawKey := r.LocalKeyExtractor(&parents[i])
		key := normalizeKey(rawKey)
		if key != nil && !mapLocalKeyToIsSeen[key] {
			mapLocalKeyToIsSeen[key] = true
			localKeys = append(localKeys, key)
		}
	}

	if len(localKeys) == 0 {
		for i := range parents {
			r.Hydrator(&parents[i], nil)
		}
		return nil
	}

	pivotTable := r.PivotTable
	if r.PivotSchema != "" && !strings.Contains(r.PivotTable, ".") {
		pivotTable = fmt.Sprintf("%s.%s", r.PivotSchema, r.PivotTable)
	}

	sqlPivot := fmt.Sprintf(
		"SELECT %s, %s FROM %s WHERE %s = ANY($1)",
		r.PivotForeignKeyToParent,
		r.PivotForeignKeyToChild,
		pivotTable,
		r.PivotForeignKeyToParent,
	)

	rows, err := db.Query(ctxEagerLoad, sqlPivot, localKeys)
	if err != nil {
		return fmt.Errorf("failed querying pivot table -> %w", err)
	}
	defer rows.Close()

	mapParentKeyToChildKeys := make(map[any][]any)
	mapChildKeyToIsSeen := make(map[any]bool)
	var childKeys []any

	for rows.Next() {
		var parentKey any
		var childKey any
		err = rows.Scan(&parentKey, &childKey)
		if err != nil {
			return fmt.Errorf("failed scanning pivot row -> %w", err)
		}
		pKey := normalizeKey(parentKey)
		cKey := normalizeKey(childKey)
		mapParentKeyToChildKeys[pKey] = append(
			mapParentKeyToChildKeys[pKey],
			cKey,
		)
		if cKey != nil && !mapChildKeyToIsSeen[cKey] {
			mapChildKeyToIsSeen[cKey] = true
			childKeys = append(childKeys, cKey)
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("pivot row iteration error -> %w", err)
	}

	if len(childKeys) == 0 {
		for i := range parents {
			r.Hydrator(&parents[i], nil)
		}
		return nil
	}

	qb := r.ChildQueryFactory(db)
	for _, customizer := range r.Customizers {
		qb = customizer(qb)
	}

	qb.Where(where.InWhere{
		Column:  r.ChildForeignKeyColumn.ColumnName(),
		Values:  childKeys,
		Boolean: contract.BoolAnd,
	})

	children, err := qb.Get(ctxEagerLoad)
	if err != nil {
		return fmt.Errorf("failed querying relation children -> %w", err)
	}

	mapChildKeyToChild := make(map[any]ChildModel)
	for i := range children {
		cKey := normalizeKey(r.ChildKeyExtractor(&children[i]))
		mapChildKeyToChild[cKey] = children[i]
	}

	for i := range parents {
		pKey := normalizeKey(r.LocalKeyExtractor(&parents[i]))
		cKeys := mapParentKeyToChildKeys[pKey]
		parentChildren := make([]ChildModel, 0, len(cKeys))
		for _, cKey := range cKeys {
			child, exists := mapChildKeyToChild[cKey]
			if exists {
				parentChildren = append(parentChildren, child)
			}
		}
		r.Hydrator(&parents[i], parentChildren)
	}

	return nil
}

func normalizeKey(val any) any {
	val = derefValue(val)
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint:
		return int64(v)
	case uint64:
		return int64(v)
	case uint32:
		return int64(v)
	case uint16:
		return int64(v)
	case uint8:
		return int64(v)
	case [16]byte:
		return fmt.Sprintf(
			"%x-%x-%x-%x-%x",
			v[0:4],
			v[4:6],
			v[6:8],
			v[8:10],
			v[10:16],
		)
	case pgtype.UUID:
		if !v.Valid {
			return nil
		}
		return fmt.Sprintf(
			"%x-%x-%x-%x-%x",
			v.Bytes[0:4],
			v.Bytes[4:6],
			v.Bytes[6:8],
			v.Bytes[8:10],
			v.Bytes[10:16],
		)
	case string:
		return strings.ToLower(v)
	}
	return val
}
