package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		wg         sync.WaitGroup
		errCount   int32                 // Атомарный счётчик ошибок
		tasksCh    = make(chan Task)     // Канал для задач
		stopSignal = make(chan struct{}) // Канал для сигнализации об остановке
		once       sync.Once             // Для безопасного закрытия stopSignal
	)

	// Функция для безопасного закрытия канала stopSignal
	closeStopSignal := func() {
		once.Do(func() {
			close(stopSignal)
		})
	}

	// Запуск n воркеров
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return // Канал закрыт, завершение работы
					}
					// Выполнение задачи
					if err := task(); err != nil {
						// Увеличиваем счётчик ошибок
						if atomic.AddInt32(&errCount, 1) >= int32(m) {
							closeStopSignal() // Если превысили лимит ошибок — сигнализируем об остановке
							return
						}
					}
				case <-stopSignal:
					return // Получен сигнал остановки
				}
			}
		}()
	}

	// Отправка задач в воркеры
	for _, task := range tasks {
		select {
		case tasksCh <- task:
		case <-stopSignal:
			break
		}
	}
	close(tasksCh)

	// Ожидание завершения всех воркеров
	wg.Wait()

	// Проверка, превышен ли лимит ошибок
	if atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
