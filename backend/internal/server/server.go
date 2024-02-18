package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"calcflow/backend/internal/orchestrator"
	"calcflow/backend/internal/task"

	"github.com/Knetic/govaluate"
)

// Server представляет HTTP-сервер для обработки запросов.
type Server struct {
	orchestrator *orchestrator.Orchestrator
}

// NewServer создает новый экземпляр HTTP-сервера с заданным оркестратором.
func NewServer(o *orchestrator.Orchestrator) *Server {
	return &Server{
		orchestrator: o,
	}
}

// Добавление вычисление нового арифметического выражения.
func (s *Server) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {

	// Генерации TaskID на стороне хендлера
	// Что бы его вернуть в случае HTTP 500
	taskID := generateTaskID()

	// Извлечение requestID и expression из JSON-тела запроса
	decoder := json.NewDecoder(r.Body)
	var requestBody map[string]string
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestID := requestBody["id"]
	expression := requestBody["expression"]

	// Проверка валидности выражения
	if !isValidExpression(expression) {
		http.Error(w, "Invalid expression", http.StatusBadRequest)
		return
	}

	// Проверка наличия requestID в базе данных
	if requestID != "" {
		if unq, err := s.AlreadyExistsRequestID(requestID); err != nil || unq {
			http.Error(w, "Request ID already exists", http.StatusOK)
			return
		}
	} else {
		http.Error(w, "Request ID is required", http.StatusBadRequest)
		return
	}

	// Добавляем вычисление в оркестратор
	errOrch := s.orchestrator.AddCalculation(expression, taskID, requestID)
	if errOrch != nil {
		http.Error(w, errOrch.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ID добавленной задачи
	responseData := map[string]string{"task_id": taskID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

// Получение списка выражений со статусами.
func (s *Server) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем список выражений с их статусами
	expressions, err := s.orchestrator.GetExpressions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем список выражений в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expressions)
}

func (s *Server) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем requestID из query параметра
	requestID := r.URL.Query().Get("requestID")

	// Получаем выражение по его идентификатору
	task, err := s.orchestrator.GetExpressionByID(requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем выражение в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Обновление времени выполнения для каждой операции POST-запросом.
func (s *Server) UpdateOperationsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	var newRequestTime task.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&newRequestTime); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем данные в БД
	if err := s.orchestrator.UpdateCalculateTime(newRequestTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

// Получение списка доступных операций со временем их выполения.
func (s *Server) GetAvailableOperationsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем доступные операции с временем выполнения
	operations, err := s.orchestrator.GetAvailableOperations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем список операций в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operations)
}

// Получение задачи для выполения. ПОКА ХЗХЗХЗХЗХЗХЗХЗХЗХЗХЗХЗ
// func (s *Server) getTaskForExecutionHandler(w http.ResponseWriter, r *http.Request) {
// 	// Проверяем метод запроса
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Получаем задачу для выполнения
// 	task, err := s.orchestrator.GetTaskForExecution()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Отправляем задачу в формате JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(task)
// }

// Приём результата обработки данных. ПОКА ХЗХЗХЗХЗХЗХЗХЗХЗХЗХЗХЗ
// func (s *Server) receiveResultHandler(w http.ResponseWriter, r *http.Request) {
// 	// Проверяем метод запроса
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Получаем параметры запроса
// 	taskID := r.FormValue("taskID")
// 	result := r.FormValue("result")
// 	errStr := r.FormValue("error")

// 	// Преобразуем ошибку в тип error
// 	var err error
// 	if errStr != "" {
// 		err = errors.New(errStr)
// 	}

// 	// Принимаем результат обработки данных
// 	err = s.orchestrator.ReceiveResult(taskID, result, err)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Отправляем успешный ответ
// 	w.WriteHeader(http.StatusOK)
// }

// Функция для проверки валидности выражения.
func isValidExpression(expression string) bool {
	_, err := govaluate.NewEvaluableExpression(expression)
	return err == nil
}

// Функция для проверки уникальности requestID в базе данных.
func (s *Server) AlreadyExistsRequestID(requestID string) (bool, error) {
	unique, err := s.orchestrator.AlreadyExistsRequest(requestID)
	if err != nil {
		return unique, err
	}

	return unique, nil
}

// generateTaskID генерирует уникальный идентификатор для задачи.
func generateTaskID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
