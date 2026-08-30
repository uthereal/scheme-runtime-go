package contract

// AggregateFunction represents a valid SQL aggregate function name.
type AggregateFunction string

// AggregateState captures aggregate selection parameters.
type AggregateState struct {
	// Function is the aggregate function name (e.g., COUNT, SUM).
	Function AggregateFunction
	// Column is the target column name for the aggregate.
	Column string
}

// AggCount represents the COUNT aggregate function.
const AggCount AggregateFunction = "COUNT"

// AggSum represents the SUM aggregate function.
const AggSum AggregateFunction = "SUM"

// AggAvg represents the AVG aggregate function.
const AggAvg AggregateFunction = "AVG"

// AggMin represents the MIN aggregate function.
const AggMin AggregateFunction = "MIN"

// AggMax represents the MAX aggregate function.
const AggMax AggregateFunction = "MAX"
