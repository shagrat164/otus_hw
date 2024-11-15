package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открытие исходного файла
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Получение информации о файле
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Проверка на поддержку файла (длина должна быть известна)
	if !srcInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	srcSize := srcInfo.Size()

	// Проверка смещения
	if offset > srcSize {
		return ErrOffsetExceedsFileSize
	}

	// Корректировка лимита, если он превышает размер файла
	if limit == 0 || offset+limit > srcSize {
		limit = srcSize - offset
	}

	// Открытие файла назначения
	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Переход к указанному смещению в исходном файле
	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek source file: %w", err)
	}

	// Инициализация прогресс-бара
	bar := progressbar.DefaultBytes(limit, "Copying")

	// Ограниченное чтение и запись с отслеживанием прогресса
	reader := io.LimitReader(srcFile, limit)
	writer := io.MultiWriter(dstFile, bar)

	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
