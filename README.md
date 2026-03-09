# Go Test Task

тестовое задание для позиции **Go Developer**

## 📋 содержание

- [Task 1](#task-1)
- [Task 2](#task-2)
- [Task 3](#task-3)
- [установка](#установка)
- [структура проекта](#структура-проекта)

---

## Task 1

CRUD приложение для управления задачами в PostgreSQL

### функциональность:

- добавление новых задач
- просмотр всех задач
- обновление статуса задачи
- удаление задач
- автоматическое создание таблицы при запуске
- использование подготовленных выражений

### запуск:

```bash
cd task1
export DATABASE_URL="postgres://postgres:пароль@localhost:5432/go_test_task?sslmode=disable"
go run main.go
```

### пример:
```bash
=== МЕНЮ ===
1. добавить задачу
2. показать все задачи
3. обновить задачу
4. удалить задачу
5. выход
выберите: 1
название: изучить Go
описание: пройти тестовое задание
статус: new
✓ задача добавлена
```
## Task 2

архивация файлов и очистка папок

### функциональность:

- создание ZIP архива из файлов папки
- удаление всех файлов из указанной папки
- рекурсивная обработка файлов
- сохранение структуры папок в архиве

### запуск:

```bash
cd task2
go run main.go ./source_folder ./delete_folder
```

### пример:
```bash
mkdir save delete
echo "тест" > save/file1.txt
echo "удалить" > delete/old.txt
go run main.go ./save ./delete
```
### результат:

✓ файлы заархивированы в archive.zip <br> 
удален: delete\old.txt <br>
✓ файлы удалены <br>

## Task 3

проверка содержимого файла на соответствие шаблону

### функциональность:

- создание тестовых файлов
- чтение содержимого файла
- проверка на соответствие строке device=123
- удаление файла при совпадении
- сохранение файла при несовпадении

### запуск:

```bash
cd task3
go run main.go
```

### пример:
```bash
=== ТЕСТ 1: файл с device=123 ===
создан файл: test1.txt
содержимое: device=123
✓ файл test1.txt УДАЛЁН (содержит device=123)

=== ТЕСТ 2: файл без device=123 ===
создан файл: test2.txt
содержимое: device=999
✓ файл test2.txt ОСТАВЛЕН (не содержит device=123)
```
### установка

### требования:

- Go 1.21 или выше
- PostgreSQL 14+ (для Task 1)
- Git

### шаги:
```bash
git clone https://github.com/KATER-0K/go-test-task.git
cd go-test-task
go mod download
```
### структура проекта
```text
go-test-task/
├── task1/
│   ├── main.go
│   └── go.mod
├── task2/
│   ├── main.go
│   └── go.mod
├── task3/
│   ├── main.go
│   └── go.mod
├── .gitignore
├── go.mod
└── README.md
```
### автор

**Katerina Terekhina** (KATER-0K) <br>
GitHub: https://github.com/KATER-0K <br>

<div align="center">

спасибо за внимание 🙏
</div>