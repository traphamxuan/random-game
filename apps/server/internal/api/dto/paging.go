package dto

import "game-random-api/utils"

type ReqPage struct {
	Offset  *int    `form:"offset" json:"offset,omitempty"`
	Limit   *int    `form:"limit" json:"limit,omitempty"`
	OrderBy *string `form:"orderBy" json:"orderBy,omitempty"`
	SortBy  *string `form:"sortBy" json:"sortBy,omitempty"`
}

func (r ReqPage) ToPage() utils.Page {
	result := utils.Page{
		Offset: 0,
		Limit:  100,
	}
	if r.Offset != nil {
		result.Offset = *r.Offset
	}
	if r.Limit != nil {
		result.Limit = *r.Limit
	}
	if r.OrderBy != nil {
		result.OrderBy = *r.OrderBy
	}
	if r.SortBy != nil {
		result.SortBy = *r.SortBy
	}
	return result
}

type RespPage struct {
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"orderBy"`
	SortBy  string `json:"sortBy"`
	Total   int    `json:"total"`
}
