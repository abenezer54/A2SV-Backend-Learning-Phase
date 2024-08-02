package model

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

func NewTask(id int, title, description string, t time.Time) *Task {
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     t,
		Status:      "not completed",
	}
}
