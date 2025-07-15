package main

import (
	"github.com/samvibes/vexop/auth-service/app"
	"github.com/samvibes/vexop/auth-service/internal/routes"
)

func main() {
	container := app.InitApp()

	router := routes.InitRoutes(container)

	router.Run(":9000")
}
