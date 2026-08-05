package handler

import (
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

func GetTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todos []models.Todo
		var err error
		todos, err = repository.GetTodos(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
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

		todo, err := repository.GetTodoByID(pool, id)

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
				"message": err.Error(),
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
	return func(c *gin.Context) {
		var idString string = c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid id type",
			})
			return
		}

		var todoStruct UpdateTodoInput
		if err := c.ShouldBindJSON(&todoStruct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		existing, err := repository.GetTodoByID(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		title := existing.Title
		if todoStruct.Title != nil {
			title = *todoStruct.Title
		}

		if todoStruct.Title == nil && todoStruct.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "At least one of title or completed must be provided",
			})
			return
		}

		completed := existing.Completed
		if todoStruct.Completed != nil {
			completed = *todoStruct.Completed
		}

		var todo *models.Todo
		todo, err = repository.UpdateTodoByID(pool, id, title, completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": todo,
		})

	}
}

func DeleteTodoByID(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid id route parameter",
			})
			return
		}

		err = repository.DeleteTodoByID(pool, id)
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
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Todo deleted successfully!",
		})
	}
}
