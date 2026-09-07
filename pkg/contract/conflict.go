package contract

// OnConflictClause represents the compiled metadata for conflict resolution.
type OnConflictClause struct {
	Action          OnConflictAction
	ConflictColumns []string
	UpdateColumns   []string
}

// OnConflictAction represents the resolution action when an insert conflicts.
type OnConflictAction int

// OnConflictNone represents the default action of ON CONFLICT.
const OnConflictNone OnConflictAction = 0

// OnConflictDoNothing performs no operation (DO NOTHING) ON CONFLICT.
const OnConflictDoNothing OnConflictAction = 1

// OnConflictDoUpdate performs an update (DO UPDATE) ON CONFLICT.
const OnConflictDoUpdate OnConflictAction = 2
