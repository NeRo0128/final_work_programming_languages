
# Gestor de Tareas Multiparadigma en Go

## Descripción
Este proyecto es una aplicación de línea de comandos (CLI) escrita en Go que permite gestionar tareas. Entre sus funcionalidades se incluyen:

- **Crear** nuevas tareas.
- **Editar** la descripción de tareas existentes.
- **Marcar** tareas como completadas.
- **Eliminar** tareas.
- **Listar** tareas pendientes, completadas o todas.
- **Guardado automático** en un archivo JSON sin bloquear el flujo principal.

La aplicación demuestra la utilización de tres paradigmas de programación:

1. **Imperativo**: Control de flujo paso a paso, asignaciones y manipulación directa de estructuras.
2. **Funcional**: Uso de funciones puras para filtrar y transformar listas de tareas sin mutar datos originales.
3. **Asincrónico**: Concurrencia en Go (goroutines y canales) para el guardado periódico en segundo plano.

---

## Tabla de Contenidos

1. [Requisitos](#requisitos)  
2. [Instalación](#instalación)  
3. [Uso](#uso)  
4. [Estructura del Proyecto](#estructura-del-proyecto)  
5. [Paradigmas de Programación](#paradigmas-de-programación)  
6. [Contribuciones](#contribuciones)  
7. [Licencia](#licencia)  

---

## Requisitos

- Go 1.18 o superior instalado en tu sistema.  
- Git (para clonar el repositorio).

---

## Instalación

1. **Clonar el repositorio**  
   ```bash
   git clone https://github.com/tu-usuario/gestor-tareas-multiparadigma-go.git
   cd gestor-tareas-multiparadigma-go

 Reemplaza `tu-usuario` por tu usuario real en GitHub.

2. **Inicializar el módulo de Go** (solo la primera vez)

   ```bash
   go mod init final_work_programming_languages
   go mod tidy
   ```

   Esto generará `go.mod` y `go.sum`. Aunque no hay dependencias externas, se recomienda mantener el módulo limpio.

3. **Compilar la aplicación**

    * Para compilar un ejecutable local:

      ```bash
      go build -o gestor-tareas main.go task.go storage.go
      ```

      Esto generará un archivo binario llamado `gestor-tareas` (o `gestor-tareas.exe` en Windows).
    * Alternativamente, puedes usar el script de ejemplo:

      ```bash
      chmod +x build.sh
      ./build.sh
      ```

      El script compilará automáticamente versiones para Linux, Windows y macOS (64-bit) dentro de la carpeta `bin/`.

---

## Uso

Una vez compilado, puedes arrancar la aplicación desde tu terminal:

```bash
./gestor-tareas
```

Verás un prompt como este:

```
Gestor de Tareas - Comandos disponibles:
  add <descripción>       : Agregar nueva tarea
  list [all|pending|done] : Listar tareas (por defecto pending)
  done <id>               : Marcar tarea como completada
  del <id>                : Eliminar tarea
  edit <id> <descripción> : Editar descripción de tarea
  help                    : Mostrar ayuda
  exit                    : Guardar y salir
```

### Comandos principales

* **Agregar nueva tarea**

  ```
  add Comprar leche
  ```

  → Crea una tarea con descripción “Comprar leche” y le asigna un ID incremental (ej. `[1]`).

* **Listar tareas pendientes (por defecto)**

  ```
  list
  ```

  → Muestra todas las tareas que aún no están completadas.

* **Listar todas las tareas**

  ```
  list all
  ```

  → Muestra pendientes y completadas (pendientes marcadas con `[ ]`, completadas con `[✔]`).

* **Listar solo tareas completadas**

  ```
  list done
  ```

  → Muestra únicamente las tareas que ya están marcadas como completadas.

* **Marcar tarea como completada**

  ```
  done 2
  ```

  → Cambia el estado de la tarea con ID 2 a “completada”.

* **Eliminar tarea**

  ```
  del 1
  ```

  → Borra la tarea con ID 1 del listado.

* **Editar descripción de una tarea**

  ```
  edit 3 Leer capítulo de “El Principito”
  ```

  → Cambia la descripción de la tarea con ID 3 a “Leer capítulo de ‘El Principito’”.

* **Mostrar ayuda**

  ```
  help
  ```

  → Vuelve a mostrar el listado de comandos disponibles.

* **Salir de la aplicación**

  ```
  exit
  ```

  → Realiza un último guardado de `tasks.json` y cierra la aplicación.

---

## Estructura del Proyecto

```
gestor-tareas-multiparadigma-go/
├── README.md
├── main.go         # Punto de entrada: parseo de comandos e interacción con el usuario.
├── task.go         # Definición de la estructura Task y funciones puramente funcionales (FilterTasks, MapDescriptions).
├── storage.go      # Persistencia en JSON (LoadTasks, SaveTasks) y goroutine de guardado automático.
└── tasks.json      # Archivo de datos que se crea la primera vez que ejecutas “./gestor-tareas”.
```

* **main.go**
  Contiene la función `main`, lee la entrada del usuario, gestiona el slice de tareas y envía señales al goroutine de guardado.

* **task.go**
  Define:

    * `type Task struct { ID int; Description string; Done bool; CreatedAt time.Time }`
    * `NextID(tasks []Task) int`
    * `ParseID(s string) int`
    * `FilterTasks(tasks []Task, predicate func(Task) bool) []Task`
    * `MapDescriptions(tasks []Task, mapper func(string) string) []string`

* **storage.go**
  Define:

    * `LoadTasks(filename string) ([]Task, error)`
    * `SaveTasks(tasks []Task)`
    * `AutoSaveLoop(tasksPtr *[]Task, saveInterval time.Duration, stopChan chan bool)`

---

## Paradigmas de Programación

### 1. Imperativo

* **Descripción**: Control de flujo paso a paso y mutación directa de estructuras.
* **Ejemplos**:

    * En `main.go`, el bucle `for` que lee comandos y las asignaciones a `tasks`, `tasks[i].Done = true`, `tasks = append(tasks, newTask)`.
    * En `storage.go`, abrir/crear archivos (`os.Create`), codificar con JSON (`encoder.Encode(tasks)`), cada línea ejecuta una instrucción que modifica el estado.

### 2. Funcional

* **Descripción**: Funciones puras que reciben datos y devuelven nuevos datos sin mutarlos.
* **Ejemplos** en `task.go`:

    * `FilterTasks(tasks []Task, predicate func(Task) bool) []Task`
      Filtra el slice original y retorna uno nuevo según la condición que define `predicate`.
    * `MapDescriptions(tasks []Task, mapper func(string) string) []string`
      Devuelve un slice de cadenas transformadas por la función `mapper`, sin tocar el slice original de `Task`.

### 3. Asincrónico (Concurrencia en Go)

* **Descripción**: Uso de goroutines y canales para ejecutar tareas en paralelo sin bloquear el flujo principal.
* **Ejemplos** en `storage.go`:

  ```go
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
  ```

    * Se crea un goroutine desde `main.go`:

      ```go
      go AutoSaveLoop(&tasks, 30*time.Second, stopChan)
      ```
    * Cada 30 segundos o al recibir señal a través de `stopChan`, se llama a `SaveTasks` para persistir el estado actual de `tasks` en `tasks.json`.