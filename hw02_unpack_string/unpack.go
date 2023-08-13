package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/example/hello/reverse"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 { // пустая строка
		return "", nil
	}

	if unicode.IsDigit(rune(str[0])) { // первый символ цифра
		return "", ErrInvalidString
	}

	var (
		strBuilder strings.Builder // для собирания строки
		count      = 1             // счётчик повтора символа
		err        error
	)
	strReverse := Reverse(str)     // обратный порядок символов исходной строки
	tmpRune := rune(strReverse[0]) // начальный символ

	for _, r := range strReverse { // перебор всех рун в строке
		if unicode.IsDigit(r) { // если цифра
			if unicode.IsDigit(tmpRune) { // если предыдущая тоже цифра
				return "", ErrInvalidString // выдать ошибку
			}
			count, err = strconv.Atoi(string(r)) // конвертнуть в int
			if err != nil {
				return "", ErrInvalidString
			}
		} else { // если это символ
			strBuilder.WriteString(strings.Repeat(string(r), count)) // добавить к строке повторив count раз
			count = 1                                                // сбросить счётчик повторов
		}
		tmpRune = r // записать текущий символ
	}
	return Reverse(strBuilder.String()), nil
}

func Reverse(str string) string {
	return reverse.String(str)
}
