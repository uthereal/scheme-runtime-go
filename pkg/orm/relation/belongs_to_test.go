package relation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/where"
)

// Test_BelongsTo_Where tests the Where method of the BelongsTo struct.
func Test_BelongsTo_Where(t *testing.T) {
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

	rel := BelongsTo[parent, child, mutator]{
		QueryBuilder: qb,
	}

	var dummyFilter contract.Where
	dummyFilter = where.BasicWhere{Column: "dummy"}
	clonedRel := rel.Where(dummyFilter)

	assert.NotNil(t, clonedRel.GetRelationQuery())
	wheres := clonedRel.GetRelationQuery().GetWheres()
	assert.Len(t, wheres, 1)
	assert.Equal(t, dummyFilter, wheres[0])
}
