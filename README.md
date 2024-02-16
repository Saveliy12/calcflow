# calcflow
CalcFlow - это приложение для обработки арифметических выражений и управления вычислительными задачами. Он включает в себя серверную часть, которая предоставляет API для добавления новых выражений, получения списка выражений, получения результатов по идентификатору выражения и других операций. В основе проекта лежит оркестратор, который управляет задачами и их выполнением, и агенты, которые фактически выполняют вычисления. В проекте используются технологии Golang, Gorilla Mux для маршрутизации HTTP запросов, а также GORM для работы с базой данных SQLite.

## Примеры работы:
### 1. Добавление вычисления арифметического выражения

**URL**: `/add-calculation`

**Метод**: `POST`

**Параметры запроса**:

- `id`: Уникальный идентификатор запроса
- `expression`: Арифметическое выражение для вычисления

**Примеры curl-запросов**:

bashCopy code

`# Пример 1: Добавление вычисления с новым уникальным идентификатором curl -X POST -H "Content-Type: application/json" -d '{"id": "unique_request_id_1", "expression": "2 + 2"}' http://localhost:8080/add-calculation  # Пример 2: Добавление вычисления с существующим идентификатором (возврат HTTP 200) curl -X POST -H "Content-Type: application/json" -d '{"id": "unique_request_id_2", "expression": "3 * 4"}' http://localhost:8080/add-calculation`

### 2. Получение списка выражений со статусами

**URL**: `/get-expressions`

**Метод**: `GET`

**Пример curl-запроса**:

bashCopy code

`curl http://localhost:8080/get-expressions`

### 3. Получение значения выражения по его идентификатору

**URL**: `/get-expression`

**Метод**: `GET`

**Параметры запроса**:

- `requestID`: Уникальный идентификатор запроса

**Пример curl-запроса**:

bashCopy code

`curl http://localhost:8080/get-expression?requestID=unique_request_id_1`

### 4. Получение списка доступных операций со временем их выполнения

**URL**: `/get-available-operations`

**Метод**: `GET`

**Пример curl-запроса**:

bashCopy code

`curl http://localhost:8080/get-available-operations`

### 5. Обновление времени выполнения для каждой операции

**URL**: `/update-operations`

**Метод**: `POST`

**Параметры запроса**: JSON-объект с полями `Summation`, `Subtraction`, `Multiplication`, `Division`, представляющими время выполнения для каждой операции.

**Пример curl-запроса**:

bashCopy code

`curl -X POST -H "Content-Type: application/json" -d '{"Summation": 10, "Subtraction": 15, "Multiplication": 20, "Division": 25}' http://localhost:8080/update-operations`

### 6. Получение задачи для выполнения (ПОКА НЕ РЕАЛИЗОВАНО)

**URL**: `/get-task-for-execution`

**Метод**: `GET`

**Пример curl-запроса**:

bashCopy code

`curl http://localhost:8080/get-task-for-execution`

### 7. Приём результата обработки данных (ПОКА НЕ РЕАЛИЗОВАНО)

**URL**: `/receive-result`

**Метод**: `POST`

**Параметры запроса**:

- `taskID`: Идентификатор задачи
- `result`: Результат выполнения задачи
- `error`: Ошибка (если есть)

**Пример curl-запроса**:

bashCopy code

`curl -X POST -d "taskID=unique_task_id&result=42&error=" http://localhost:8080/receive-result`
