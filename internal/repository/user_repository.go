package repository

import (
	"context"
	"time"

	"github.com/Promise111/go-rest-api-todo/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *models.Users) (*models.Users, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	var query string = `
	INSERT INTO users (email, password, username) 
	VALUES($1, $2, $3)
	RETURNING id, email, username, created_at, updated_at;
	`
	err = pool.QueryRow(ctx, query, user.Email, user.Password, user.Username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.Users, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	var query string = `
	SELECT id, email, password, username, created_at, updated_at 
	FROM users 
	WHERE email = $1
	`
	var user models.Users
	err = pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.Users, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	var query string = `
	SELECT id, email, password, username, created_at, updated_at 
	FROM users 
	WHERE id = $1
	`
	var user models.Users
	err = pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(pool *pgxpool.Pool, username string) (*models.Users, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	var query string = `
	SELECT id, email, password, username, created_at, updated_at 
	FROM users 
	WHERE username = $1
	`

	var user models.Users
	if err = pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
