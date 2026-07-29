package main

import (
	"fmt"

	"github.com/Promise111/go-rest-api-todo/internal/router"
)

func main() {
	fmt.Println("🚀 Todo API Server started!")
	fmt.Println("💹 Project setup complete")
	fmt.Println("📁 Folder structure all set")
	r := router.Router()

	r.Run(":8080")
}
