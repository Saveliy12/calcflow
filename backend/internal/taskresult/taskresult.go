package taskresult

import "calcflow/backend/internal/task"

// ResultProcessor интерфейс для обработки результатов выполнения задач.
type ResultProcessor interface {
	ReceiveResult(task *task.Task) error
	GetAvailableOperations() (*task.CalculationRequest, error)
	EnqueueTask(task *task.Task)
}

type TaskProcessor interface {
	EnqueueTask(task *task.Task)
}
