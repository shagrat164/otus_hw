package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/example/hello/reverse"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 { // пустая строка
		return "", nil
	}

	firstRune, _ := utf8.DecodeRuneInString(str) // получить первый символ

	if unicode.IsDigit(firstRune) { // первый символ цифра
		return "", ErrInvalidString
	}

	var (
		strBuilder strings.Builder // для собирания строки
		count      = 1             // счётчик повтора символа
		tmpRune    rune
	)
	strReverse := reverse.String(str) // обратный порядок символов исходной строки

	for len(strReverse) > 0 {
		curRune, sizeRune := utf8.DecodeRuneInString(strReverse) // получить первый символ строки и его размер в байтах
		if curRune == utf8.RuneError {
			return "", ErrInvalidString // выдать ошибку если сивол некорректен
		}

		if unicode.IsDigit(curRune) { // если цифра
			if unicode.IsDigit(tmpRune) { // если предыдущая тоже цифра
				return "", ErrInvalidString // выдать ошибку
			}
			count, _ = strconv.Atoi(string(curRune)) // конвертнуть в int
		} else { // если это символ
			strBuilder.WriteString(strings.Repeat(string(curRune), count)) // добавить к строке повторив count раз
			count = 1                                                      // сбросить счётчик повторов
		}
		tmpRune = curRune                  // сохранить текущий символ
		strReverse = strReverse[sizeRune:] // обрезать массив
	}

	return reverse.String(strBuilder.String()), nil
}
