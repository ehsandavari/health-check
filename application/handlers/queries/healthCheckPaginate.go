package queries

import (
	"health-check/application/common"
)

type SHealthCheckPaginateQuery struct {
	paginateQuery common.PaginateQuery
}

func NewHealthCheckPaginateQuery(paginateQuery common.PaginateQuery) SHealthCheckPaginateQuery {
	return SHealthCheckPaginateQuery{
		paginateQuery: paginateQuery,
	}
}
