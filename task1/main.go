package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
}

// для JSON файла
type TaskFromFile struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var db *sql.DB
var scanner *bufio.Scanner

func main() {
	// подключение к БД
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres123@localhost:5432/go_test_task?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// проверка подключения
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// создаем таблицу
	createTable()

	// создаем сканер для чтения ввода
	scanner = bufio.NewScanner(os.Stdin)

	// меню
	for {
		showMenu()
		choice := readLine()
		
		switch choice {
		case "1":
			createTask()
		case "2":
			readTasks()
		case "3":
			updateTask()
		case "4":
			deleteTask()
		case "5":
			importTasksFromFile()
		case "6":
			fmt.Println("выход из программы...")
			return
		default:
			fmt.Println("неверный выбор. попробуйте снова.")
		}
	}
}

func readLine() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("таблица готова!")
}

func showMenu() {
	fmt.Println("\n=== МЕНЮ ===")
	fmt.Println("1. добавить задачу")
	fmt.Println("2. показать все задачи")
	fmt.Println("3. обновить задачу")
	fmt.Println("4. удалить задачу")
	fmt.Println("5. импортировать задачи из файлов")
	fmt.Println("6. выход")
	fmt.Print("выберите: ")
}

func createTask() {
	fmt.Print("название: ")
	title := readLine()
	
	fmt.Print("описание: ")
	description := readLine()
	
	fmt.Print("статус: ")
	status := readLine()

	_, err := db.Exec("INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)",
		title, description, status)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ задача добавлена")
}

func readTasks() {
	rows, err := db.Query("SELECT id, title, description, status, created_at FROM tasks ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\n=== ЗАДАЧИ ===")
	fmt.Printf("%-3s | %-20s | %-30s | %-10s | %s\n", "ID", "название", "описание", "статус", "дата")
	fmt.Println(strings.Repeat("-", 90))
	
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-3d | %-20s | %-30s | %-10s | %s\n", 
			t.ID, t.Title, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func updateTask() {
	fmt.Print("ID задачи: ")
	idStr := readLine()
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("неверный ID")
		return
	}
	
	fmt.Print("новый статус: ")
	status := readLine()

	result, err := db.Exec("UPDATE tasks SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		log.Fatal(err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("задача не найдена")
	} else {
		fmt.Println("✓ задача обновлена")
	}
}

func deleteTask() {
	fmt.Print("ID задачи: ")
	idStr := readLine()
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("неверный ID")
		return
	}

	result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("задача не найдена")
	} else {
		fmt.Println("✓ задача удалена")
	}
}

// импорт задач из JSON файлов
func importTasksFromFile() {
	fmt.Print("введите путь к папке с файлами (например ./files): ")
	filesDir := readLine()

	// ищем все JSON файлы в папке
	files, err := filepath.Glob(filepath.Join(filesDir, "*.json"))
	if err != nil {
		fmt.Printf("ошибка поиска файлов: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("файлы не найдены в папке")
		return
	}

	fmt.Printf("найдено файлов: %d\n", len(files))

	// создаём канал для результатов
	results := make(chan string, len(files))

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	errorCount := 0

	for result := range results {
		fmt.Println(result)
		if strings.HasPrefix(result, "✓") {
			successCount++
		} else {
			errorCount++
		}
	}

	fmt.Println("\n=== ИТОГИ ИМПОРТА ===")
	fmt.Printf("успешно: %d\n", successCount)
	fmt.Printf("ошибки: %d\n", errorCount)
}

func processFile(filePath string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	data, err := os.ReadFile(filePath)
	if err != nil {
		results <- fmt.Sprintf("ошибка чтения %s: %v", filePath, err)
		return
	}

	// парсим JSON
	var task TaskFromFile
	err = json.Unmarshal(data, &task)
	if err != nil {
		results <- fmt.Sprintf("ошибка парсинга %s: %v", filePath, err)
		return
	}

	// записываем в базу данных
	_, err = db.Exec(
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)",
		task.Title, task.Description, task.Status,
	)
	if err != nil {
		results <- fmt.Sprintf("ошибка записи %s: %v", filePath, err)
		return
	}

	results <- fmt.Sprintf("✓ успешно: %s", filePath)
}