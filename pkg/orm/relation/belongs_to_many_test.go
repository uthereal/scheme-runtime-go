package relation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_BelongsToMany_Where tests the Where method of the BelongsToMany struct.
func Test_BelongsToMany_Where(t *testing.T) {
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

	rel := BelongsToMany[parent, child, mutator]{
		ChildQueryFactory: func(
			_ contract.DB,
		) *orm.QueryBuilder[child, mutator] {
			return qb.Clone()
		},
	}

	var dummyFilter contract.Where
	dummyFilter = where.BasicWhere{Column: "dummy"}
	clonedRel := rel.Where(dummyFilter)

	assert.NotNil(t, clonedRel.GetRelationQuery())
	wheres := clonedRel.GetRelationQuery().GetWheres()
	assert.Len(t, wheres, 1)
	assert.Equal(t, dummyFilter, wheres[0])
}

// Test_BelongsToMany_Customise_Isolation tests that customizing maintains
// immutability.
func Test_BelongsToMany_Customise_Isolation(t *testing.T) {
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

	sharedRel := BelongsToMany[parent, child, mutator]{
		ChildQueryFactory: func(
			_ contract.DB,
		) *orm.QueryBuilder[child, mutator] {
			return qb.Clone()
		},
	}

	filterA := where.BasicWhere{Column: "type", Value: "A"}
	relA := sharedRel.Where(filterA)

	filterB := where.BasicWhere{Column: "type", Value: "B"}
	relB := sharedRel.Where(filterB)

	assert.Len(t, sharedRel.Customizers, 0)
	assert.Len(t, sharedRel.GetRelationQuery().GetWheres(), 0)

	wheresA := relA.GetRelationQuery().GetWheres()
	assert.Len(t, wheresA, 1)
	assert.Equal(t, filterA, wheresA[0])

	wheresB := relB.GetRelationQuery().GetWheres()
	assert.Len(t, wheresB, 1)
	assert.Equal(t, filterB, wheresB[0])
}

// Test_BelongsToMany_EagerLoad_Empty tests empty parent slices.
func Test_BelongsToMany_EagerLoad_Empty(t *testing.T) {
	type parent struct{}
	type child struct{}
	type mutator struct{}

	rel := BelongsToMany[parent, child, mutator]{}
	err := rel.EagerLoad(nil, nil, nil)
	assert.NoError(t, err)
}
