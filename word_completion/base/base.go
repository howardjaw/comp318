package base

import (
	"sort"
	"strings"
)

type Completer struct {
	words []string
}

func New(dict map[string]int) Completer {
	words := make([]string, 0, len(dict))

	for word := range dict {
		words = append(words, word)
	}

	sort.Strings(words)

	return Completer{
		words: words,
	}
}

func (c Completer) Complete(prefix string) []string {
	results := []string{}

	for _, word := range c.words {
		if strings.HasPrefix(word, prefix) {
			results = append(results, word)
		}

		if len(results) == 10 {
			break
		}
	}
	return results
}
