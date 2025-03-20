package main

import (
	"bingo/middleware"
	"bingo/routes"
	"bingo/validation"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4020"
	}

	validation.InitValidation()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CorsMiddleware)
	routes.InitRoutes(r)

	fmt.Println("List of Routes")
	for _, route := range r.Routes() {
		green := "\033[32m"
		reset := "\033[0m"
		fmt.Printf("%s%s%s %s\n", green, route.Method, reset, route.Path)
	}

	fmt.Println("server running on port:", port)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
