package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"calcflow/backend/internal/task"
)

type Store struct {
	db *gorm.DB
}

// Создание сущности базы данных
func New(path string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't open database: %v", err)
	}

	// Выполните миграцию таблицы, если это необходимо
	err = db.AutoMigrate(&task.Task{})
	if err != nil {
		return nil, fmt.Errorf("can't migrate database: %v", err)
	}

	return &Store{db: db}, nil
}

// Закрытие соединения с базой данных
func (s *Store) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Создание необходимых таблиц
func (s *Store) CreateTables() error {
	// Создание таблицы Tasks
	err := s.db.AutoMigrate(&task.Task{})
	if err != nil {
		return err
	}

	// Создание таблицы CalculationRequest
	err = s.db.AutoMigrate(&task.CalculationRequest{})
	if err != nil {
		return err
	}

	return nil
}

// Добавление новой задачи в таблицу `Tasks`
func (s *Store) NewTask(task *task.Task) error {
	result := s.db.Create(task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Получение задачи по requestID из таблицы `Tasks`
func (s *Store) GetTaskByID(requestID string) (*task.Task, error) {
	var task task.Task
	result := s.db.Where("request_id = ?", requestID).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// Получение всех задач из таблицы `Tasks`
func (s *Store) GetAllTasks() ([]*task.Task, error) {
	var tasks []*task.Task
	result := s.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// Проверка на существование выражения с таким requestID
func (s *Store) AlreadyExistsRequest(requestID string) (bool, error) {
	var count int64
	result := s.db.Model(&task.Task{}).Where("request_id = ?", requestID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Получение времени выполнения каждой операции из таблицы `CalculationRequest`
func (s *Store) GetCalculateTime() (*task.CalculationRequest, error) {
	var calcRequest task.CalculationRequest
	result := s.db.First(&calcRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	return &calcRequest, nil
}

// Обновление значений времени выполнения для каждой операции
func (s *Store) UpdateCalculateTime(request task.CalculationRequest) error {
	result := s.db.Save(&request)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Обновление данных задачи в таблице `Tasks` после того, как выражение будет посчитано
func (s *Store) UpdateTask(task *task.Task) error {
	result := s.db.Save(task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
