package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(txt string) []string {
	set := strings.Fields(txt)

	dict := make(map[string]int)
	for _, word := range set {
		dict[word]++
	}

	wordsSlice := make([]string, 0, len(dict)) // слайс размером с мапу dict

	for key := range dict {
		wordsSlice = append(wordsSlice, key) // заполнить слайс уникальными словами
	}

	// сортировка Девид Блейн стайл
	// по слову берётся количество повторений в dict
	// и сортируется по убыванию или по алфавиту, если
	// количество повторов одинаково
	sort.Slice(wordsSlice, func(i, j int) bool {
		a, b := wordsSlice[i], wordsSlice[j]
		if dict[a] == dict[b] {
			return a < b
		}
		return dict[a] > dict[b]
	})

	if len(wordsSlice) > 10 {
		wordsSlice = wordsSlice[:10] // максимум 10 слов на выходе
	}

	return wordsSlice
}
