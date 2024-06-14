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
