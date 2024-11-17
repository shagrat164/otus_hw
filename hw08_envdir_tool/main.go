package main

import (
	"log"
	"os"
)

func main() {
	// Проверка что верно переданы аргументы
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s /path/to/env/dir command arg1 arg2 ...", os.Args[0])
	}

	envDir := os.Args[1]
	command := os.Args[2:]
	env, err := ReadDir(envDir) // Чтение окружения из директории
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	exitCode := RunCmd(command, env)
	os.Exit(exitCode)
}
