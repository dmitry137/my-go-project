package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dmitry137/my-go-project/handlers"
	"github.com/dmitry137/my-go-project/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

func main() {
	// Получаем параметры подключения из переменных окружения
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDbName := os.Getenv("PG_DBNAME")
	pgSSLMode := os.Getenv("PG_SSLMODE")

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgDbName,
		pgSSLMode,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT,
            status TEXT CHECK (status IN ('new', 'in_progress', 'done')) NOT NULL DEFAULT 'new',
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP NOT NULL DEFAULT NOW()
        )`); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	storage := storage.NewPostgresStorage(conn)
	app := fiber.New()

	app.Post("/tasks", handlers.CreateTaskHandler(storage))
	app.Get("/tasks", handlers.GetTasksHandler(storage))
	app.Put("/tasks/:id", handlers.UpdateTaskHandler(storage))
	app.Delete("/tasks/:id", handlers.DeleteTaskHandler(storage))

	log.Fatal(app.Listen(":3000"))
}
