package handler

import (
	"log/slog"
	"net/http"

	"github.com/Promise111/go-rest-api-todo/internal/models"
	"github.com/Promise111/go-rest-api-todo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoInput CreateTodoInput
		var err error
		err = c.ShouldBindJSON(&todoInput)
		if err != nil {
			slog.Error("Something went wrong", "message", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		var todo *models.Todo
		todo, err = repository.CreateTodo(pool, todoInput.Title, todoInput.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  true,
			"message": "Todo created successfully",
			"data":    todo,
		})
	}
}
