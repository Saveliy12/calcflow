package main

import (
	"calcflow/backend/internal/agent"
	"calcflow/backend/internal/database"
	"calcflow/backend/internal/orchestrator"
	"calcflow/backend/internal/server"
	"calcflow/backend/internal/task"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Реализация интерфейса TaskProcessor
type MyProcessor struct {
	// здесь могут быть поля, необходимые для обработки результатов задач
}

func (p *MyProcessor) EnqueueTask(task *task.Task) {
	// Здесь можно выполнить необходимые действия перед добавлением задачи в очередь агента
}

func (p *MyProcessor) ReceiveResult(task *task.Task) error {
	// Здесь можно выполнить необходимые действия с результатами выполнения задачи
	return nil
}

// Функция для создания экземпляра MyProcessor
func NewMyProcessor() *MyProcessor {
	return &MyProcessor{}
}

func main() {
	// Инициализация базы данных
	db, err := database.New("database.db")
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Создание необходимых таблиц
	// db.CreateTables()

	// Создание объекта MyProcessor
	processor := NewMyProcessor()

	// Создание оркестратора
	orchestrator, err := orchestrator.NewOrchestrator(db, processor)
	if err != nil {
		log.Fatalf("Ошибка при создании оркестратора: %v", err)
	}

	// Создание агента (или агентов)
	agent.NewAgent("AgentName", 100, orchestrator)

	// Инициализация и запуск сервера
	s := server.NewServer(orchestrator)
	router := mux.NewRouter()

	// Обработчики запросов
	router.HandleFunc("/add-calculation", s.AddExpressionHandler).Methods("POST")
	router.HandleFunc("/get-expressions", s.GetExpressionsHandler).Methods("GET")
	router.HandleFunc("/get-expression", s.GetExpressionByIDHandler).Methods("GET")
	router.HandleFunc("/update-operations", s.UpdateOperationsHandler).Methods("POST")
	router.HandleFunc("/get-available-operations", s.GetAvailableOperationsHandler).Methods("GET")

	// Запуск сервера

	fmt.Println("Сервер запущен на :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
