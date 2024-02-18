package task

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Task представляет структуру арифметического выражения.
type Task struct {
	ID         string        `json:"id"`
	RequestID  string        `json:"X-Request-id"`
	Expression string        `json:"expression"`
	Status     string        `json:"status"`
	Result     string        `json:"result"`
	Created    time.Time     `json:"created"`
	Finished   time.Time     `json:"finished"`
	Duration   time.Duration `json:"duration"`
}

// CalculationRequest представляет значения выполнения каждой арифметической операции.
type CalculationRequest struct {
	Summation      string `json:"summation"`
	Subtraction    string `json:"subtraction"`
	Multiplication string `json:"multiplication"`
	Division       string `json:"division"`
}

// Value реализует интерфейс database/sql/driver.Valuer для CalcRequest.
func (cr CalculationRequest) Value() (driver.Value, error) {
	return json.Marshal(cr)
}

// Scan реализует интерфейс database/sql/driver.Scanner для CalcRequest.
func (cr *CalculationRequest) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan: не удалось преобразовать в []byte")
	}
	return json.Unmarshal(bytes, &cr)
}
