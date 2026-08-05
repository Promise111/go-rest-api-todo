package repository

import (
	"context"
	"time"

	"github.com/Promise111/go-rest-api-todo/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var query string = `
	INSERT INTO todos (title, completed) 
	VALUES ($1, $2)
	RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at 
	FROM todos 
	ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos = []models.Todo{}

	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at 
	FROM todos 
	WHERE id = $1
	`

	var todo models.Todo
	var err error = pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoByID(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var query string = `
	UPDATE todos 
	SET title = $1, completed = $2 updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo
	var err error = pool.QueryRow(ctx, query, id, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodoByID(pool *pgxpool.Pool, id int) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	DELETE FROM todos 
	WHERE id = $1
	`
	res, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
