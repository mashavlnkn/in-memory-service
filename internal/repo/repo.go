package repo

import (
	"context"
	"sync"
	"time"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

type repository struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

// Repository - интерфейс с методом создания задачи
type Repository interface {
	CreateTask(ctx context.Context, task Task) (int, error)
	GetTaskByID(ctx context.Context, id int) (*Task, error)
	GetTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, id int, updated Task) error
	DeleteTask(ctx context.Context, id int) error // Создание задачи
}

func NewRepository() Repository {
	return &repository{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}
func (r *repository) CreateTask(ctx context.Context, task Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = "pending" // Статус по умолчанию

	r.tasks[task.ID] = task
	r.nextID++

	return task.ID, nil
}
func (r *repository) GetTasks(ctx context.Context) ([]Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var tasks []Task
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *repository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}
	return &task, nil
}
func (r *repository) DeleteTask(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return nil
	}
	delete(r.tasks, id)
	return nil
}
func (r *repository) UpdateTask(ctx context.Context, id int, updated Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return nil
	}

	updated.ID = id
	updated.UpdatedAt = time.Now()
	r.tasks[id] = updated
	return nil
}

// CreateTask - вставка новой задачи в таблицу tasks
