package handlers

import (
	"strconv"

	"github.com/dmitry137/my-go-project/models"
	"github.com/dmitry137/my-go-project/storage"
	"github.com/gofiber/fiber/v2"
)

func isValidStatus(status string) bool {
	switch status {
	case "new", "in_progress", "done":
		return true
	default:
		return false
	}
}

func CreateTaskHandler(s storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task models.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if task.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
		}

		if task.Status == "" {
			task.Status = "new"
		}

		if !isValidStatus(task.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
		}

		createdTask, err := s.CreateTask(task)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
		}

		return c.Status(fiber.StatusCreated).JSON(createdTask)
	}
}

func GetTasksHandler(s storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tasks, err := s.GetTasks()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
		}
		return c.JSON(tasks)
	}
}

func UpdateTaskHandler(s storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		var task models.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if task.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
		}

		if !isValidStatus(task.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
		}

		task.ID = id

		existingTask, err := s.GetTaskByID(id)
		if err != nil {
			if err.Error() == "task not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch task"})
		}

		if task.Description == nil {
			task.Description = existingTask.Description
		}

		if err := s.UpdateTask(task); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
		}

		updatedTask, err := s.GetTaskByID(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated task"})
		}

		return c.JSON(updatedTask)
	}
}

func DeleteTaskHandler(s storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		if err := s.DeleteTask(id); err != nil {
			if err.Error() == "task not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
