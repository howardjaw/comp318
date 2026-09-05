package frequency

import (
	"sort"
	"strings"
)

type wordFreq struct {
	word  string
	count int
}
type Completer struct {
	words []wordFreq
}

func New(dict map[string]int) Completer {
	words := make([]wordFreq, 0, len(dict))

	for word, count := range dict {
		words = append(words, wordFreq{
			word:  word,
			count: count,
		})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count
	})

	return Completer{
		words: words,
	}
}

func (c Completer) Complete(prefix string) []string {
	results := []string{}

	for _, item := range c.words {
		if strings.HasPrefix(item.word, prefix) {
			results = append(results, item.word)
		}

		if len(results) == 10 {
			break
		}
	}

	return results
}
