package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Promise111/go-rest-api-todo/internal/config"
	"github.com/Promise111/go-rest-api-todo/internal/models"
	"github.com/Promise111/go-rest-api-todo/internal/repository"
	"github.com/Promise111/go-rest-api-todo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required_without=Username,omitempty,email"`
	Username string `json:"username" binding:"required_without=Email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func RegisterUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest RegisterRequest
		var err error

		if err = c.ShouldBindJSON(&registerRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Password must be at least 6 characters long",
			})
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		var user *models.Users = &models.Users{
			Email:    registerRequest.Email,
			Password: string(hashedPass),
			Username: registerRequest.Username,
		}

		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Email or Username already registered",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong!",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  true,
			"message": "User created successfully!",
			"data":    createdUser,
		})
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		var err error
		if err = c.ShouldBindJSON(&loginRequest); err != nil {
			slog.Error("Error", "err", err.Error())
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				fieldErrors := make(map[string]string, len(ve))
				for _, fe := range ve {
					fieldErrors[utils.JsonFieldNameLogin(fe)] = utils.ValidationMessageLogin(fe)
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Validation failed",
					"errors":  fieldErrors,
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid JSON",
			})
			return
		}
		
	}
}
