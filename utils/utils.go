package utils

import "sort"

// Convert map → slice for sorting
type kv struct {
	Message string
	Count   int
}

func SortMap(m map[string]int) []string {
	var list []kv
	for msg, count := range m {
		list = append(list, kv{msg, count})
	}

	// Sort by count (descending)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	limit := 5
	if len(list) < 5 {
		limit = len(list)
	}

	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = list[i].Message
	}
	return result
}
