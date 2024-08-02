package service

import (
	"errors"

	"task-manager-api/model"
)

type TaskManager struct {
	Tasks map[int]model.Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make(map[int]model.Task),
	}
}

func (tm *TaskManager) GetAllTasks() []model.Task {
	tasks := []model.Task{}
	for _, task := range tm.Tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (tm *TaskManager) GetTask(taskID int) (model.Task, error) {
	task, ok := tm.Tasks[taskID]
	if !ok {
		return model.Task{}, errors.New("couldn't find task")
	}
	return task, nil
}

func (tm *TaskManager) CreateTask(task model.Task) error {
	_, ok := tm.Tasks[task.ID]
	if ok {
		return errors.New("task already exists")
	}
	tm.Tasks[task.ID] = task
	return nil
}

func (tm *TaskManager) UpdateTask(taskID int, t model.Task) error {
	_, ok := tm.Tasks[taskID]
	if !ok {
		return errors.New("couldn't find task")
	}

	if t.Title == "" || t.Description == "" || t.DueDate.IsZero() || t.Status == "" {
		return errors.New("all fields must be filled to update with PUT request")
	}

	if newID := t.ID; newID != taskID {
		if _, ok := tm.Tasks[newID]; ok {
			return errors.New("ID you provided already exists")
		}
	}

	tm.Tasks[taskID] = t

	return nil
}

func (tm *TaskManager) PartialUpdateTask(taskID int, t model.Task) error {
	_, ok := tm.Tasks[taskID]
	if !ok {
		return errors.New("couldn't find task")
	}

	if t.ID != taskID {
		return errors.New("cannot update id with patch request")
	}

	tm.Tasks[taskID] = t

	return nil
}

func (tm *TaskManager) DeleteTask(taskID int) error {
	_, ok := tm.Tasks[taskID]
	if !ok {
		return errors.New("couldn't find task")
	}
	delete(tm.Tasks, taskID)
	return nil
}
