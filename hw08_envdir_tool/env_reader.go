package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidFileName = errors.New("invalid file name containing '='")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue // пропустить директории
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			return nil, ErrInvalidFileName
		}

		filePath := filepath.Join(dir, name)
		content, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer content.Close()

		scanner := bufio.NewScanner(content)
		value := ""
		if scanner.Scan() {
			value = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Убрать пробелы и табуляцию в конце строки
		value = strings.TrimRight(value, " \t")

		if value == "" && !scanner.Scan() {
			// Если файл пустой
			env[name] = EnvValue{
				Value:      "",
				NeedRemove: true, // Пометить переменную для удаления
			}
		} else {
			// Заменить терминальные нули
			value = strings.ReplaceAll(value, "\x00", "\n")
			env[name] = EnvValue{
				Value:      value,
				NeedRemove: false,
			}
		}
	}
	return env, nil
}
