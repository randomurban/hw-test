package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	r := regexp.MustCompile(`[\t\n\f\r ,.!:;(){}\[\]*_']+`)
	tokens := r.Split(input, -1)
	counts := map[string]int{}

	for _, t := range tokens {
		if t != "-" && t != "" {
			counts[strings.ToLower(t)]++
		}
	}

	keys := make([]string, 0)
	for k := range counts {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if counts[keys[i]] == counts[keys[j]] {
			return keys[i] < keys[j]
		}
		return counts[keys[i]] > counts[keys[j]]
	})
	if len(keys) > 10 {
		return keys[:10]
	}
	return keys
}
