package apiservice

type PageInfo struct {
	HasPreviousPage bool
	HasNextPage     bool
	StartCursor     *Cursor
	EndCursor       *Cursor
}
