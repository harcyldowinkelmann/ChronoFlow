package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Task struct {
	ID          int    `json:"id"`
	ScheduleID  int    `json:"scheduleId"`
	RowIndex    int    `json:"rowIndex"`
	DayOfWeek   int    `json:"dayOfWeek"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Description string `json:"description"`
	IsCompleted bool   `json:"isCompleted"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Tags        []Tag  `json:"tags"`
}

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	appDir := filepath.Join(configDir, "ChronoFlow")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatalf("Failed to create application directory: %v", err)
	}

	dbPath := filepath.Join(appDir, "cronograma_v3.db")

	// Initialize SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	a.db = db

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			schedule_id INTEGER DEFAULT 1,
			row_index INTEGER,
			day_of_week INTEGER,
			start_time TEXT,
			end_time TEXT,
			description TEXT,
			is_completed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			color TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS task_tags (
			task_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (task_id, tag_id),
			FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
		);`,
	}

	for _, q := range queries {
		if _, err := a.db.Exec(q); err != nil {
			log.Fatalf("Failed to execute query: %v", err)
		}
	}
}

// GetTasks returns all tasks from the database with their tags
func (a *App) GetTasks() ([]Task, error) {
	rows, err := a.db.Query("SELECT id, schedule_id, row_index, day_of_week, start_time, end_time, description, is_completed, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ScheduleID, &t.RowIndex, &t.DayOfWeek, &t.StartTime, &t.EndTime, &t.Description, &t.IsCompleted, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		
		// Fetch tags for this task
		tagRows, err := a.db.Query("SELECT t.id, t.name, t.color FROM tags t INNER JOIN task_tags tt ON t.id = tt.tag_id WHERE tt.task_id = ?", t.ID)
		if err == nil {
			var tags []Tag
			for tagRows.Next() {
				var tag Tag
				if err := tagRows.Scan(&tag.ID, &tag.Name, &tag.Color); err == nil {
					tags = append(tags, tag)
				}
			}
			t.Tags = tags
			tagRows.Close()
		}

		if t.Tags == nil {
			t.Tags = []Tag{}
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// SaveTask creates a new task or updates an existing one
func (a *App) SaveTask(task Task) (Task, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	var taskId int

	// Check if exists based on row_index and day_of_week
	err := a.db.QueryRow("SELECT id FROM tasks WHERE row_index = ? AND day_of_week = ?", task.RowIndex, task.DayOfWeek).Scan(&taskId)
	
	if err == sql.ErrNoRows {
		// Insert
		res, err := a.db.Exec(`
			INSERT INTO tasks (schedule_id, row_index, day_of_week, start_time, end_time, description, is_completed, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			1, task.RowIndex, task.DayOfWeek, task.StartTime, task.EndTime, task.Description, task.IsCompleted, now, now)
		if err != nil {
			return Task{}, fmt.Errorf("failed to insert task: %v", err)
		}
		id, _ := res.LastInsertId()
		taskId = int(id)
	} else if err == nil {
		// Update
		_, err = a.db.Exec(`
			UPDATE tasks SET start_time = ?, end_time = ?, description = ?, is_completed = ?, updated_at = ? 
			WHERE id = ?`,
			task.StartTime, task.EndTime, task.Description, task.IsCompleted, now, taskId)
		if err != nil {
			return Task{}, fmt.Errorf("failed to update task: %v", err)
		}
	} else {
		return Task{}, fmt.Errorf("failed to check task: %v", err)
	}

	// Update tags
	_, _ = a.db.Exec("DELETE FROM task_tags WHERE task_id = ?", taskId)
	for _, tag := range task.Tags {
		_, _ = a.db.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", taskId, tag.ID)
	}

	// Return updated task
	task.ID = taskId
	task.UpdatedAt = now
	if task.CreatedAt == "" {
		task.CreatedAt = now
	}
	return task, nil
}

// DeleteTask removes a task by ID
func (a *App) DeleteTask(id int) error {
	_, err := a.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// GetTags returns all tags
func (a *App) GetTags() ([]Tag, error) {
	rows, err := a.db.Query("SELECT id, name, color FROM tags")
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %v", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %v", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// SaveTag creates a new tag
func (a *App) SaveTag(name string, color string) (Tag, error) {
	res, err := a.db.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", name, color)
	if err != nil {
		return Tag{}, fmt.Errorf("failed to insert tag: %v", err)
	}
	id, _ := res.LastInsertId()
	return Tag{ID: int(id), Name: name, Color: color}, nil
}

// DeleteTag removes a tag
func (a *App) DeleteTag(id int) error {
	_, err := a.db.Exec("DELETE FROM tags WHERE id = ?", id)
	return err
}
