package dto

func IncludeRelationIfValid(relations *[]string, query *bool, relation string) {
	if query != nil && *query {
		*relations = append(*relations, relation)
	}
}
