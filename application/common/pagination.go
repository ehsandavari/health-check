package common

import (
	"math"
)

type PaginateResult[T any] struct {
	Page       uint
	PerPage    uint
	TotalPage  uint
	TotalItems uint64
	Items      []T
}

func NewPaginateResult[T any](items []T, page uint, perPage uint, totalItems uint64) *PaginateResult[T] {
	paginateResult := &PaginateResult[T]{Items: items, Page: page, PerPage: perPage, TotalItems: totalItems}

	paginateResult.TotalPage = getTotalPages(totalItems, perPage)

	return paginateResult
}

func getTotalPages(totalCount uint64, perPage uint) uint {
	d := float64(totalCount) / float64(perPage)
	return uint(math.Ceil(d))
}

type ComparisonType string

const (
	ComparisonTypeIn       ComparisonType = "IN (?)"
	ComparisonTypeEquals   ComparisonType = "= ?"
	ComparisonTypeContains ComparisonType = "LIKE ?"
)

func (r ComparisonType) String() string {
	return string(r)
}

func (r ComparisonType) IsValid() bool {
	switch r {
	case ComparisonTypeIn,
		ComparisonTypeEquals,
		ComparisonTypeContains:
		return true
	default:
		return false
	}
}

type FilterQuery struct {
	Key        string         `binding:"required"`
	Comparison ComparisonType `binding:"required,enum"`
	Value      string         `binding:"required"`
}

type PaginateQuery struct {
	Page    uint `binding:"required" default:"1"`
	PerPage uint `binding:"required" default:"10"`
	OrderBy string
	Filters []FilterQuery `binding:"dive"`
}

func (q *PaginateQuery) GetOffset() uint {
	return (q.Page - 1) * q.PerPage
}

func (q *PaginateQuery) GetLimit() uint {
	return q.PerPage
}

func (q *PaginateQuery) GetOrderBy() string {
	return q.OrderBy
}

func (q *PaginateQuery) GetPage() uint {
	return q.Page
}

func (q *PaginateQuery) GetPerPage() uint {
	return q.PerPage
}
