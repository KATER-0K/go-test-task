package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("использование: task2.exe <папка_для_архива> <папка_для_удаления>")
		return
	}

	saveDir := os.Args[1]
	deleteDir := os.Args[2]

	// архивируем
	err := zipDirectory(saveDir, "archive.zip")
	if err != nil {
		fmt.Printf("ошибка архивации: %v\n", err)
	} else {
		fmt.Println("✓ файлы заархивированы в archive.zip")
	}

	// удаляем
	err = deleteFiles(deleteDir)
	if err != nil {
		fmt.Printf("ошибка удаления: %v\n", err)
	} else {
		fmt.Println("✓ файлы удалены")
	}
}

func zipDirectory(sourceDir, zipFile string) error {
	outFile, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		zipEntry, _ := zipWriter.Create(relPath)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})
}

func deleteFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			os.Remove(path)
			fmt.Printf("удален: %s\n", path)
		}
		return nil
	})
}