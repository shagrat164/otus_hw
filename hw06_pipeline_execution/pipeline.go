package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		res := make(Bi)
		close(res)
		return res
	}

	outCh := in
	for _, stage := range stages {
		outCh = func(in In) (out Out) {
			bridgeCh := make(Bi) // Создаётся промежуточный канал для связи стадий
			go func() {
				defer close(bridgeCh) // Закрытие канала по завершению работы

				for {
					select {
					case <-done: // Прекратить работу при закрытии `done`
						<-in
						return
					case data, ok := <-in:
						if !ok { // Завершение если вход закрыт
							return
						}
						select {
						case <-done: // Проверка `done` перед отправкой
						case bridgeCh <- data:
						}
					}
				}
			}()
			return stage(bridgeCh)
		}(outCh)
	}
	return outCh
}
