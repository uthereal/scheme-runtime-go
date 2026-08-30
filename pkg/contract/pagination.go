package contract

// PaginateResult holds items and whether there are more records.
type PaginateResult[Model any] struct {
	// Items is the slice of models returned for the current page.
	Items []Model

	// HasMore indicates whether there are additional pages available.
	HasMore bool
}

// PaginateAip158Result holds the page items and the next page token
// for subsequent requests.
type PaginateAip158Result[Model any] struct {
	// Items is the slice of models returned for the current page.
	Items []Model

	// NextPageToken is the token used to retrieve the next page.
	NextPageToken string
}

// PaginateCountResult holds items and the total matching record count.
type PaginateCountResult[Model any] struct {
	// Items is the slice of models returned for the current page.
	Items []Model

	// TotalCount is the total number of records matching the query.
	TotalCount int64
}
