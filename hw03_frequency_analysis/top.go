package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {

	tokens := strings.Fields(input)
	counts := map[string]int{}

	for _, t := range tokens {
		counts[t]++

	}

	var keys []string
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sort.SliceStable(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})

	var res []string
	for i, k := range keys {
		if i > 9 {
			break
		}
		res = append(res, k)
	}
	return res
}
