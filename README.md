# my-go-project
API для управления задачами с использованием Go, Fiber и PostgreSQL.


## Запуск через Docker

1. Склонируйте репозиторий:
- git clone https://github.com/dmitry137/my-go-project
- cd my-go-project
2. Запустите проект:
- docker compose up --build

API будет доступно на http://localhost:3000 (упрощенно в рамках задания)

## Примеры запросов
Создать задачу
- curl -X POST -H "Content-Type: application/json" -d '{"title":"My Task"}' http://localhost:3000/tasks
Получить все задачи
- curl http://localhost:3000/tasks