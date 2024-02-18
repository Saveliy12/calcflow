package agent

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"calcflow/backend/internal/task"
	"calcflow/backend/internal/taskresult"

	"github.com/Knetic/govaluate"
)

// Agent представляет вычислительный агент.
type Agent struct {
	Name      string          // Имя агента
	WorkQueue chan *task.Task // Канал-очередь, откуда агент будет брать задачи
	mu        sync.Mutex
	processor taskresult.ResultProcessor
}

// NewAgent создает новый экземпляр агента.
func NewAgent(name string, workQueueSize int, processor taskresult.ResultProcessor) *Agent {
	return &Agent{
		Name:      name,
		WorkQueue: make(chan *task.Task, workQueueSize),
		processor: processor,
	}
}

// Start запускает агента и начинает обработку задач в его очереди.
func (a *Agent) Start() {
	for {
		var task *task.Task

		// Блокировка мьютекса для доступа к каналу
		a.mu.Lock()

		// Получение задачи из очереди
		select {
		case task = <-a.WorkQueue:
		default:
			// Если очередь пуста, освобождаем мьютекс и ждем
			a.mu.Unlock()
			time.Sleep(time.Millisecond * 100) // Даем шанс другим горутинам
			continue
		}

		// Освобождение мьютекса после получения задачи
		a.mu.Unlock()

		// Обработка задачи
		go a.processTask(task)
	}
}

// ExecuteExpression выполняет вычисление арифметического выражения.
func (a *Agent) ExecuteExpression(expressionStr string, calcRequest task.CalculationRequest) (string, error) {
	// Создаем новое выражение
	expression, err := govaluate.NewEvaluableExpression(expressionStr)
	if err != nil {
		return "", err
	}

	// Подсчитываем количество операций в выражении
	opCounts := map[string]int{
		"summation":      strings.Count(expressionStr, "+"),
		"subtraction":    strings.Count(expressionStr, "-"),
		"multiplication": strings.Count(expressionStr, "*"),
		"division":       strings.Count(expressionStr, "/"),
	}

	summation, _ := time.ParseDuration(calcRequest.Summation)
	subtraction, _ := time.ParseDuration(calcRequest.Subtraction)
	multiplication, _ := time.ParseDuration(calcRequest.Multiplication)
	division, _ := time.ParseDuration(calcRequest.Division)

	// Рассчитываем общее время выполнения выражения
	totalDuration := time.Duration(0)
	for op, count := range opCounts {
		switch op {
		case "summation":
			totalDuration += summation * time.Duration(count)
		case "subtraction":
			totalDuration += subtraction * time.Duration(count)
		case "multiplication":
			totalDuration += multiplication * time.Duration(count)
		case "division":
			totalDuration += division * time.Duration(count)
		}
	}

	// Создаем таймер с общим временем выполнения
	timer := time.NewTimer(totalDuration)

	// Запускаем горутину для выполнения выражения
	resultChan := make(chan string)
	go func() {
		defer close(resultChan)
		result, err := expression.Evaluate(nil)
		if err != nil {
			resultChan <- ""
			return
		}
		resultChan <- fmt.Sprintf("%v", result)
	}()

	// Ждем результат или таймаут
	select {
	case result := <-resultChan:
		// Если получен результат из канала, возвращаем его
		return result, nil
	case <-timer.C:
		// Если истекло время таймаута, возвращаем ошибку
		return "", errors.New("превышено время выполнения выражения")
	}
}

// processTask обрабатывает задачу и отправляет результат обратно оркестратору.
// processTask обрабатывает задачу и отправляет результат обратно оркестратору.
func (a *Agent) processTask(taskToWork *task.Task) {
	var maxAttempts = 3
	var retryDelay = time.Millisecond * 100

	// Получаем время выполнения для каждой операции от оркестратора из таблицы CalculationRequest
	for attempts := 0; attempts < maxAttempts; attempts++ {
		calcRequest, err := a.processor.GetAvailableOperations()
		if err != nil {
			log.Printf("Ошибка получения времени выполнения операций: %v. Повторная попытка через %v", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Проверяем, все ли значения в CalculationRequest заполнены
		if isEmpty(calcRequest) {
			// Если хотя бы одно значение пусто, устанавливаем время выполнения каждой операции в ноль
			calcRequest = &task.CalculationRequest{
				Summation:      "0s",
				Subtraction:    "0s",
				Multiplication: "0s",
				Division:       "0s",
			}
		}

		// Обработка задачи
		result, err := a.ExecuteExpression(taskToWork.Expression, *calcRequest)
		if err != nil || result == "" {
			taskToWork.Status = "error" // Меняем статус вычисления выражения на "error"
			taskToWork.Result = ""
		} else {
			taskToWork.Status = "completed" // Меняем статус вычисления выражения на "completed"
			taskToWork.Result = result
		}

		// Отправка результата обратно оркестратору
		a.processor.ReceiveResult(taskToWork)
		return
	}

	// Если не удалось выполнить задачу после нескольких попыток, устанавливаем статус "error"
	taskToWork.Status = "error"
	taskToWork.Result = ""
	a.processor.ReceiveResult(taskToWork)
}

// isEmpty проверяет, является ли структура CalculationRequest пустой (все поля пусты или равны нулю)
func isEmpty(calcRequest *task.CalculationRequest) bool {
	return calcRequest.Summation == "" && calcRequest.Subtraction == "" && calcRequest.Multiplication == "" && calcRequest.Division == ""
}

// EnqueueTask добавляет задачу в очередь агента для выполнения.
func (a *Agent) EnqueueTask(task *task.Task) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.WorkQueue <- task
}
