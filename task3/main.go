package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// тест 1: файл с device=123
	fmt.Println("=== ТЕСТ 1: Файл с device=123 ===")
	testFile1 := "test1.txt"
	os.WriteFile(testFile1, []byte("device=123"), 0644)
	fmt.Printf("Создан файл: %s\n", testFile1)

	shouldDelete := checkFile(testFile1)
	if shouldDelete {
		os.Remove(testFile1)
		fmt.Printf("✓ Файл %s УДАЛЕН (содержит device=123)\n", testFile1)
	} else {
		fmt.Printf("Файл %s оставлен\n", testFile1)
	}

	// тест 2: файл без device=123
	fmt.Println("\n=== ТЕСТ 2: Файл без device=123 ===")
	testFile2 := "test2.txt"
	os.WriteFile(testFile2, []byte("device=999"), 0644)
	fmt.Printf("Создан файл: %s\n", testFile2)

	shouldDelete = checkFile(testFile2)
	if shouldDelete {
		os.Remove(testFile2)
		fmt.Printf("Файл %s УДАЛЕН\n", testFile2)
	} else {
		fmt.Printf("✓ Файл %s ОСТАВЛЕН (не содержит device=123)\n", testFile2)
	}
}

func checkFile(filename string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Ошибка чтения: %v\n", err)
		return false
	}

	content := strings.TrimSpace(string(data))
	fmt.Printf("Содержимое: %s\n", content)

	return content == "device=123"
}