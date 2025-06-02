package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func printHelp() {
	fmt.Println("Gestor de Tareas - Comandos disponibles:")
	fmt.Println("  add <descripción>       : Agregar nueva tarea")
	fmt.Println("  list [all|pending|done] : Listar tareas (por defecto pending)")
	fmt.Println("  done <id>               : Marcar tarea como completada")
	fmt.Println("  del <id>                : Eliminar tarea")
	fmt.Println("  edit <id> <descripción> : Editar descripción de tarea")
	fmt.Println("  help                    : Mostrar ayuda")
	fmt.Println("  exit                    : Guardar y salir")
}

func main() {
	// Cargar tareas iniciales
	tasks, err := LoadTasks("tasks.json")
	if err != nil {
		fmt.Println("No se encontró file tasks.json. Se creará uno nuevo al guardar.")
		tasks = []Task{}
	}

	// Canal para sincronizar guardado
	stopChan := make(chan bool)
	// Iniciar goroutine de guardado automático cada 30 segundos
	go AutoSaveLoop(&tasks, 30*time.Second, stopChan)

	scanner := bufio.NewScanner(os.Stdin)
	printHelp()
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := parts[0]
		switch cmd {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Uso: add <descripción>")
				continue
			}
			desc := strings.Join(parts[1:], " ")
			newTask := Task{
				ID:          NextID(tasks),
				Description: desc,
				Done:        false,
				CreatedAt:   time.Now(),
			}
			tasks = append(tasks, newTask)
			fmt.Printf("Tarea [%d] agregada.\n", newTask.ID)
			// Avisar al goroutine para guardar inmediatamente
			stopChan <- false

		case "list":
			filter := "pending"
			if len(parts) >= 2 {
				filter = parts[1]
			}
			var filtered []Task
			switch filter {
			case "all":
				filtered = tasks
			case "done":
				filtered = FilterTasks(tasks, func(t Task) bool { return t.Done })
			default: // "pending"
				filtered = FilterTasks(tasks, func(t Task) bool { return !t.Done })
			}
			if len(filtered) == 0 {
				fmt.Println("No hay tareas para mostrar.")
				continue
			}
			for _, t := range filtered {
				status := "[ ]"
				if t.Done {
					status = "[✔]"
				}
				fmt.Printf("  %s %d: %s\n", status, t.ID, t.Description)
			}

		case "done":
			if len(parts) != 2 {
				fmt.Println("Uso: done <id>")
				continue
			}
			id := ParseID(parts[1])
			updated := false
			for i, t := range tasks {
				if t.ID == id {
					tasks[i].Done = true
					fmt.Printf("Tarea [%d] marcada como completada.\n", id)
					updated = true
					break
				}
			}
			if !updated {
				fmt.Printf("No se encontró tarea con ID %d.\n", id)
			} else {
				stopChan <- false
			}

		case "del":
			if len(parts) != 2 {
				fmt.Println("Uso: del <id>")
				continue
			}
			id := ParseID(parts[1])
			newTasks := []Task{}
			deleted := false
			for _, t := range tasks {
				if t.ID == id {
					deleted = true
					continue
				}
				newTasks = append(newTasks, t)
			}
			if deleted {
				tasks = newTasks
				fmt.Printf("Tarea [%d] eliminada.\n", id)
				stopChan <- false
			} else {
				fmt.Printf("No se encontró tarea con ID %d.\n", id)
			}

		case "edit":
			if len(parts) < 3 {
				fmt.Println("Uso: edit <id> <descripción>")
				continue
			}
			id := ParseID(parts[1])
			newDesc := strings.Join(parts[2:], " ")
			edited := false
			for i, t := range tasks {
				if t.ID == id {
					tasks[i].Description = newDesc
					fmt.Printf("Tarea [%d] actualizada.\n", id)
					edited = true
					break
				}
			}
			if !edited {
				fmt.Printf("No se encontró tarea con ID %d.\n", id)
			} else {
				stopChan <- false
			}

		case "help":
			printHelp()

		case "exit":
			// Señalar goroutine para guardado final y terminar
			stopChan <- true
			fmt.Println("Guardando y saliendo…")
			time.Sleep(100 * time.Millisecond) // Esperar breve para guardar
			return

		default:
			fmt.Println("Comando no reconocido. Escribe 'help' para ver opciones.")
		}
	}
}
