package relation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_Relation_GetRelationQuery tests the GetRelationQuery method.
func Test_Relation_GetRelationQuery(t *testing.T) {
	type parent struct{}
	type child struct{}
	type mutator struct{}

	rel := Relation[parent, child, mutator]{
		QueryBuilder: nil,
	}

	assert.Nil(t, rel.GetRelationQuery())
}

// Test_Relation_EagerLoad_Empty parents tests the EagerLoad method with empty.
func Test_Relation_EagerLoad_Empty(t *testing.T) {
	type parent struct{}
	type child struct{}
	type mutator struct{}

	rel := Relation[parent, child, mutator]{}
	err := rel.EagerLoad(context.Background(), nil, nil)
	assert.NoError(t, err)
}

// Test_Relation_Constrain_Multiple tests that multiple constraints
// (limits, offsets, wheres) can be fluent-chained onto a relation
// and correctly applied when resolved.
func Test_Relation_Constrain_Multiple(t *testing.T) {
	type parent struct{}
	type child struct{}
	type mutator struct{}

	qb := orm.NewQueryBuilder[child, mutator](
		nil,
		nil,
		contract.TableMetadata[child]{},
		nil,
		nil,
	)

	rel := HasMany[parent, child, mutator]{
		Relation: Relation[parent, child, mutator]{
			QueryBuilder: qb,
		},
	}

	dummyFilter := where.BasicWhere{Column: "active"}

	// Fluent chain both Where and Constrain
	customizedRel := rel.
		Where(dummyFilter).
		Constrain(func(
			q *orm.QueryBuilder[child, mutator],
		) *orm.QueryBuilder[child, mutator] {
			return q.Limit(10).Offset(5)
		})

	// Get compiled QueryStateProvider
	qState := customizedRel.GetRelationQuery()
	assert.NotNil(t, qState)

	// Verify wheres
	wheres := qState.GetWheres()
	assert.Len(t, wheres, 1)
	assert.Equal(t, dummyFilter, wheres[0])

	// Verify limit
	limit, ok := qState.GetLimit()
	assert.True(t, ok)
	assert.Equal(t, uint64(10), limit)

	// Verify offset
	offset, ok := qState.GetOffset()
	assert.True(t, ok)
	assert.Equal(t, uint64(5), offset)
}

// Test_Relation_Customise_Isolation tests that customizing a relation
// maintains strict immutability, ensuring concurrent/separate queries
// do not leak query criteria into the shared original relation instance.
func Test_Relation_Customise_Isolation(t *testing.T) {
	type parent struct{}
	type child struct{}
	type mutator struct{}

	qb := orm.NewQueryBuilder[child, mutator](
		nil,
		nil,
		contract.TableMetadata[child]{},
		nil,
		nil,
	)

	// This is our simulated shared / global relationship instance (e.g. from
	// Schema structure)
	sharedRel := HasMany[parent, child, mutator]{
		Relation: Relation[parent, child, mutator]{
			QueryBuilder: qb,
		},
	}

	// First query customization
	filterA := where.BasicWhere{Column: "type", Value: "A"}
	relA := sharedRel.Where(filterA)

	// Second query customization
	filterB := where.BasicWhere{Column: "type", Value: "B"}
	relB := sharedRel.Where(filterB)

	// Verify sharedRel remains completely empty/uncustomized
	assert.Len(t, sharedRel.Customizers, 0)
	assert.Len(t, sharedRel.GetRelationQuery().GetWheres(), 0)

	// Verify relA only has filterA
	wheresA := relA.GetRelationQuery().GetWheres()
	assert.Len(t, wheresA, 1)
	assert.Equal(t, filterA, wheresA[0])

	// Verify relB only has filterB (no leakage from filterA!)
	wheresB := relB.GetRelationQuery().GetWheres()
	assert.Len(t, wheresB, 1)
	assert.Equal(t, filterB, wheresB[0])
}
