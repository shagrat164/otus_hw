package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	// Проверить, что команда безопасна
	executablePath, err := exec.LookPath(cmd[0]) // Полный путь к команде
	if err != nil {
		// Если команда не найдена
		return 127
	}

	// Создать команду
	command := exec.Command(executablePath, cmd[1:]...) // Безопасный вызов
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Установить окружение
	environ := os.Environ()
	for key, val := range env {
		if val.NeedRemove {
			// Удалить переменную среды
			for i, e := range environ {
				if len(e) > len(key) && e[:len(key)+1] == key+"=" {
					environ = append(environ[:i], environ[i+1:]...)
					break
				}
			}
		} else {
			// Добавить новую переменную среды
			environ = append(environ, key+"="+val.Value)
		}
	}
	command.Env = environ

	// Запуск команды
	if err := command.Run(); err != nil {
		// Обработка ошибок при запуске
		if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
