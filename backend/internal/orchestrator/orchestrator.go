package orchestrator

import (
	"sync"
	"time"

	"calcflow/backend/internal/database"
	"calcflow/backend/internal/task"
	"calcflow/backend/internal/taskresult"
)

// Orchestrator представляет оркестратор, управляющий задачами.
type Orchestrator struct {
	tasks     map[string]*task.Task // Мапа для хранения задач
	mu        sync.Mutex
	processor taskresult.TaskProcessor
	db        *database.Store // Ссылка на сущность базы данных
}

// NewOrchestrator создает новый экземпляр оркестратора.
func NewOrchestrator(db *database.Store, processor taskresult.TaskProcessor) (*Orchestrator, error) {
	return &Orchestrator{
		tasks:     make(map[string]*task.Task),
		db:        db,
		processor: processor,
	}, nil
}

// AddCalculation добавляет новое арифметическое выражение для вычисления
func (o *Orchestrator) AddCalculation(expression, taskID, requestID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	task := &task.Task{
		ID:         taskID,
		RequestID:  requestID,
		Expression: expression,
		Status:     "pending",
		Created:    time.Now(),
	}

	// Сохранение задачи в базе данных
	err := o.db.NewTask(task)
	if err != nil {
		return err
	}

	// Отправка задачи на выполнение агенту
	o.processor.EnqueueTask(task)

	// Возвращаем ID задачи
	return nil
}

// GetExpressionByID возвращает значение арифметического выражения по его идентификатору
func (o *Orchestrator) GetExpressionByID(requestID string) (*task.Task, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	task, err := o.db.GetTaskByID(requestID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetExpressionsWithStatus возвращает список арифметических выражений с их статусами
func (o *Orchestrator) GetExpressions() ([]*task.Task, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	expressions, err := o.db.GetAllTasks()
	if err != nil {
		return nil, err
	}

	return expressions, nil
}

// GetAvailableOperations возвращает список доступных операций и времени их выполнения
func (o *Orchestrator) GetAvailableOperations() (*task.CalculationRequest, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	calcRequest, err := o.db.GetCalculateTime()
	if err != nil {
		return nil, err
	}

	return calcRequest, nil

}

// UpdateCalculateTim обновляет значений времени выполнения для каждой операции
func (o *Orchestrator) UpdateCalculateTime(newRequestTime task.CalculationRequest) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	err := o.db.UpdateCalculateTime(newRequestTime)
	if err != nil {
		return err
	}

	return nil

}

// ReceiveResult принимает результат обработки данных от агента
func (o *Orchestrator) ReceiveResult(task *task.Task) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	task.Finished = time.Now()               // Время окончания вычисления операции
	task.Duration = time.Since(task.Created) // Время вычисления выражения

	err := o.db.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}

// isDuplicateRequest проверяет, что такой requestID уникальный
func (o *Orchestrator) AlreadyExistsRequest(requestID string) (bool, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	unique, err := o.db.AlreadyExistsRequest(requestID)
	if err != nil {
		return false, err
	}

	return unique, nil

}
