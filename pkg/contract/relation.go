package contract

import "context"

// Relation marks any typed relation object that can be preloaded,
// strictly bound to its ParentModel.
type Relation[ParentModel any] interface {
	// GetRelationQuery returns the underlying query provider for the relation.
	GetRelationQuery() QueryStateProvider

	// EagerLoad accepts fetched parent models and hydrates them directly.
	EagerLoad(
		ctxEagerLoad context.Context,
		db DB,
		parents []ParentModel,
	) error
}
