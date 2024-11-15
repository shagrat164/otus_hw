package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	// Создаем директорию для тестов
	testDir := t.TempDir()
	inputFile := filepath.Join(testDir, "input.txt")
	outputFile := filepath.Join(testDir, "output.txt")

	// Заполняем исходный файл тестовыми данными
	content := []byte("Hello, World!")
	err := os.WriteFile(inputFile, content, 0o644)
	assert.NoError(t, err, "Failed to create input file")

	t.Run("FullCopy", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 0)
		assert.NoError(t, err, "Full copy should succeed")

		result, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Failed to read output file")
		assert.Equal(t, content, result, "Output file content mismatch")
	})

	t.Run("PartialCopy", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 7, 5) // Копируем "World"
		assert.NoError(t, err, "Partial copy should succeed")

		result, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Failed to read output file")
		assert.Equal(t, []byte("World"), result, "Output file content mismatch")
	})

	t.Run("OffsetExceedsSize", func(t *testing.T) {
		err := Copy(inputFile, outputFile, int64(len(content)+1), 5)
		assert.ErrorIs(t, err, ErrOffsetExceedsFileSize, "Offset exceeding file size should return error")
	})

	t.Run("UnsupportedFile", func(t *testing.T) {
		// Используем `/dev/null` или подобный файл
		unsupportedFile := "/dev/null"
		err := Copy(unsupportedFile, outputFile, 0, 10)
		assert.ErrorIs(t, err, ErrUnsupportedFile, "Unsupported file type should return error")
	})
}

func TestCopyLargeFile(t *testing.T) {
	// Создаем большой файл для проверки копирования
	testDir := t.TempDir()
	inputFile := filepath.Join(testDir, "large_input.txt")
	outputFile := filepath.Join(testDir, "large_output.txt")

	// Генерируем большой файл (например, 10 MB)
	data := make([]byte, 10*1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	err := os.WriteFile(inputFile, data, 0o644)
	assert.NoError(t, err, "Failed to create large input file")

	err = Copy(inputFile, outputFile, 0, 0)
	assert.NoError(t, err, "Copying large file should succeed")

	result, err := os.ReadFile(outputFile)
	assert.NoError(t, err, "Failed to read large output file")
	assert.Equal(t, data, result, "Output file content mismatch for large file")
}
