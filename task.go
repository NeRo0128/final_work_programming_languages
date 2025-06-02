package main

import (
	"strconv"
	"time"
)

// Task representa una tarea simple
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

// NextID retorna un ID incremental basándose en el slice actual
func NextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// ParseID intenta convertir string a int para IDs; si falla, devuelve 0 (no válido)
func ParseID(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return id
}

// FilterTasks devuelve un nuevo slice con tareas que cumplen la función predicate
func FilterTasks(tasks []Task, predicate func(Task) bool) []Task {
	result := []Task{}
	for _, t := range tasks {
		if predicate(t) {
			result = append(result, t)
		}
	}
	return result
}

// MapDescriptions devuelve un slice de strings transformando cada descripción con mapper
func MapDescriptions(tasks []Task, mapper func(string) string) []string {
	result := make([]string, len(tasks))
	for i, t := range tasks {
		result[i] = mapper(t.Description)
	}
	return result
}
