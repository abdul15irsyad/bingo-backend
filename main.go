package main

import (
	"bingo/config"
	"bingo/handler"
	"bingo/lib"
	"bingo/middleware"
	"bingo/routes"
	"bingo/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	lib.InitZap()
	lib.InitValidation()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	Init(r)

	fmt.Println("List of Routes")
	for _, route := range r.Routes() {
		green := "\033[32m"
		reset := "\033[0m"
		fmt.Printf("%s%s%s %s\n", green, route.Method, reset, route.Path)
	}

	fmt.Println("server running on port:", config.Port)
	err := r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
}

func Init(r *gin.Engine) {
	dbManager := lib.NewGormManager()
	postgresDB, err := dbManager.InitPostgresDB("main", config.DBConfig)
	if err != nil {
		panic(err)
	}
	// service
	gameService := service.NewGameService()
	socketService := service.NewSocketService(2)
	userService := service.NewUserService(postgresDB)
	// middleware
	corsMiddleware := middleware.NewCorsMiddleware()
	errorMiddleware := middleware.NewErrorMiddleware()
	jwtMiddleware := middleware.NewJWTMiddleware(userService)
	// handler
	authHandler := handler.NewAuthHandler(userService)
	profileHandler := handler.NewProfileHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	socketHandler := handler.NewSocketHandler(socketService, gameService)
	// route
	rootRoute := routes.NewRootRoute()
	authRoute := routes.NewAuthRoute(authHandler)
	profileRoute := routes.NewProfileRoute(profileHandler)
	userRoute := routes.NewUserRoute(userHandler)
	socketRoute := routes.NewSocketRoute(socketHandler)

	r.Use(errorMiddleware.Handler)
	r.Use(corsMiddleware.Handler)

	rootRoute.InitRootRoute(r)
	authRoute.InitAuthRoute(r)

	r.Use(jwtMiddleware.Handler)
	profileRoute.InitProfileRoute(r)
	userRoute.InitUserRoute(r)
	socketRoute.InitSocketRoute(r)

}
