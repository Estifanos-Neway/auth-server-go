package main

import (
	"github.com/estifanos-neway/auth-server-with-go/src/handlers"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router := gin.Default()
	router.POST("/sign-up", handlers.SignUpUser)
	router.POST("/sign-in", handlers.SignInUser)

	router.Run(":8080")
}
