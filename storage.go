package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex // Mutex para evitar condiciones de carrera al acceder al slice

// LoadTasks intenta leer el archivo JSON y decodificar en []Task
func LoadTasks(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// SaveTasks guarda el slice de tareas en el archivo JSON (sobrescribe)
func SaveTasks(tasks []Task) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error al crear archivo:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		fmt.Println("Error al guardar tareas:", err)
	}
}

// AutoSaveLoop corre en un goroutine y guarda cada intervalo dado.
// Si recibe true en stopChan, hace un guardado final y retorna.
func AutoSaveLoop(tasksPtr *[]Task, saveInterval time.Duration, stopChan chan bool) {
	ticker := time.NewTicker(saveInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			SaveTasks(*tasksPtr)
		case stop := <-stopChan:
			SaveTasks(*tasksPtr)
			if stop {
				return
			}
		}
	}
}
