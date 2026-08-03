package router

import (
	handler "github.com/Promise111/go-rest-api-todo/internal/handlers"
	"github.com/Promise111/go-rest-api-todo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Router(pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	api := r.Group(utils.APIPrefix)

	{
		health := api.Group(utils.HealthPrefix)
		health.GET("", handler.HealthHandler)
	}

	{
		todos := api.Group(utils.TodosPrefix)
		todos.POST("", handler.CreateTodoHandler(pool))
	}

	return r
}
