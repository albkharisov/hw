package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type pair struct {
	word  string
	count int
}

type pairSlice []pair

func (p pairSlice) Len() int      { return len(p) }
func (p pairSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pairSlice) Less(i, j int) bool {
	var result bool = (p[i].count < p[j].count) ||
		((p[i].count == p[j].count) &&
			(p[i].word > p[j].word))
	return result
}

func (p *pairSlice) AddWord(word string) {
	found := false

	for i := range *p {
		if (*p)[i].word == word {
			found = true
			(*p)[i].count++
			break
		}
	}

	if !found {
		*p = append(*p, pair{word, 1})
	}
}

func Top10(input string) []string {
	var pl pairSlice

	for _, v := range strings.Fields(input) {
		pl.AddWord(v)
	}

	sort.Sort(sort.Reverse(pl))

	result := make([]string, 0, 10)
	for k := range pl {
		result = append(result, pl[k].word)
		if len(result) == 10 {
			break
		}
	}

	return result
}
