package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
)

var (
	re      = regexp.MustCompile(`[\s]+`)
	verbose = false
)

type mySliceType struct {
	word string // слово
	num  int    // количество
}

func Top10(txt string) []string {
	split := re.Split(txt, -1)
	set := append([]string{}, split...)

	sliceUniqueValue := getCountUniqueValue(set)

	// сортировка по переменной .num и .word, если .num одинаков
	sort.Slice(sliceUniqueValue, func(i, j int) bool {
		if sliceUniqueValue[i].num == sliceUniqueValue[j].num {
			return sliceUniqueValue[i].word < sliceUniqueValue[j].word
		}
		return sliceUniqueValue[i].num > sliceUniqueValue[j].num
	})

	sliceOnlyWord := make([]string, 0, len(sliceUniqueValue))

	if verbose {
		fmt.Printf("sliceUniqueWordAndValue: ")
		fmt.Println(sliceUniqueValue) // для наглядности
	}

	for _, element := range sliceUniqueValue { // пройтись по всему слайсу
		sliceOnlyWord = append(sliceOnlyWord, element.word) // и добавить только слова уже отсортированные
	}

	if len(sliceOnlyWord) > 10 {
		sliceOnlyWord = sliceOnlyWord[:10] // максимум 10 слов
	}

	if verbose {
		fmt.Printf("sliceOnlyWord: ")
		fmt.Println(sliceOnlyWord) // для наглядности
	}

	return sliceOnlyWord
}

func getCountUniqueValue(arr []string) []mySliceType {
	// мапа уникальных значений, где ключ это слово,
	// а значение это количество вхождений в массиве arr
	dict := make(map[string]int)
	for _, word := range arr {
		dict[word]++
	}

	resultSlice := make([]mySliceType, 0, len(dict))
	var tmpMySlice mySliceType
	// тут вытаскиваю данные в mySliceType
	for key := range dict {
		if key == "" {
			// Я не понял, почему при нескольких пробелах в map попадает строка ""
			continue // пропуск
		}
		tmpMySlice.word = key
		tmpMySlice.num = dict[key]
		resultSlice = append(resultSlice, tmpMySlice)
	}

	return resultSlice
}
