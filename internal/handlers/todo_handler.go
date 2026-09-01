package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Promise111/go-rest-api-todo/internal/models"
	"github.com/Promise111/go-rest-api-todo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}
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

		userID, ok := userIDValue.(string)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		var todo *models.Todo
		todo, err = repository.CreateTodo(pool, todoInput.Title, todoInput.Completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
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

func GetTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong",
			})
			return
		}
		userID := userIDValue.(string)
		var todos []models.Todo
		var err error
		todos, err = repository.GetTodos(pool, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Todos fetched successfully",
			"data":    todos,
			"length":  len(todos),
		})
	}
}

func GetTodoByID(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		userIDVal, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong",
			})
			return
		}

		userID := userIDVal.(string)

		todo, err := repository.GetTodoByID(pool, id, userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "Todo not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Todo fetched succesfully",
			"data":    todo,
		})
	}
}

func UpdateTodoByID(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var idParam = ctx.Param("id")
		var id int
		var err error
		
		userIDVal, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		userID := userIDVal.(string)

		id, err = strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "ID route parameter is invalid",
			})
			return
		}

		var input UpdateTodoInput
		err = ctx.ShouldBindJSON(&input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
		var existing *models.Todo
		existing, err = repository.GetTodoByID(pool, id, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "Todo not found",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}
		var title string = existing.Title
		var completed bool = existing.Completed
		if input.Title != nil {
			title = *input.Title
		}
		if input.Completed != nil {
			completed = *input.Completed
		}

		if input.Title == nil && input.Completed == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Expected at-least one of title or completed",
			})
			return
		}

		if title == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Title can not be empty",
			})
			return
		}

		var todo *models.Todo
		todo, err = repository.UpdateTodoByID(pool, id, title, completed, userID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"status":  true,
			"message": "Todo updated successfully",
			"data":    todo})
	}
}

func DeleteTodoByID(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var id int
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		userID := userIDVal.(string)

		var idParam string = c.Param("id")
		id, err = strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if err = repository.DeleteTodoByID(pool, id, userID); err != nil {
			slog.Error("Error", "error", err)
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "Todo not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		// http.StatusNoContent
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Todo deleted successfully!",
		})
	}
}
