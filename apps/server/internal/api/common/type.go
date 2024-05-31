package common

type Page struct {
	Offset  int
	Limit   int
	OrderBy string
	SortBy  string
	Total   int
}
