package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"simple-service/internal/dto"
	"simple-service/internal/repo"
	"simple-service/pkg/validator"
	"strconv"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

func (s *service) GetTasks(ctx *fiber.Ctx) error {
	tasks, err := s.repo.GetTasks(ctx.Context())
	if err != nil {
		s.log.Error("Failed to retrieve tasks", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	if len(tasks) == 0 {
		s.log.Warn("No tasks found")
		return dto.NotFoundError(ctx, "No tasks found")
	}

	response := dto.Response{
		Status: "success",
		Data:   tasks,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetTaskByID - обработчик запроса на получение задачи по ID
func (s *service) GetTaskByID(ctx *fiber.Ctx) error {
	// Получаем ID из query-параметра
	idStr := ctx.Params("id")
	if idStr == "" {
		s.log.Error("Task ID is required")
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Task ID is required")
	}

	// Преобразуем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Error("Invalid task ID format", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID format")
	}

	// Получаем задачу из репозитория
	task, err := s.repo.GetTaskByID(ctx.Context(), id)
	if err != nil {
		s.log.Error("Failed to retrieve task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Если задача не найдена, возвращаем 404
	if task == nil {
		s.log.Warn("Task not found", zap.Int("task_id", id))
		return dto.NotFoundError(ctx, "Task not found")
	}

	// Формируем JSON-ответ
	response := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// CreateTask - обработчик запроса на создание задачи
func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Вставка задачи в БД через репозиторий
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	taskID, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	// Получаем ID из параметра URL
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Error("Invalid task ID format", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID format")
	}

	// Удаляем задачу из репозитория
	err = s.repo.DeleteTask(ctx.Context(), id)
	if err != nil {
		if err.Error() == "task not found" { // ✅ 404, если задачи нет
			s.log.Warn("Task not found", zap.Int("task_id", id))
			return dto.NotFoundError(ctx, "Task not found")
		}
		s.log.Error("Failed to delete task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	// ✅ 204 No Content, если удаление успешное
	return ctx.SendStatus(fiber.StatusNoContent)
}
func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	// Получаем ID задачи из параметров URL
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Error("Invalid task ID format", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID format")
	}

	var req TaskRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Обновляем задачу в репозитории
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	err = s.repo.UpdateTask(ctx.Context(), id, task)
	if err != nil {
		if err.Error() == "task not found" { // ✅ 404, если задачи нет
			s.log.Warn("Task not found", zap.Int("task_id", id))
			return dto.NotFoundError(ctx, "Task not found")
		}
		s.log.Error("Failed to update task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": id},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
