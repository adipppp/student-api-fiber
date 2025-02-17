package main

import (
	// "fmt"
	"studentapifiber/controllers"
	// "studentapifiber/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func registerMiddlewares(app *fiber.App) error {
	// authMiddlewareHandler, err := middlewares.NewAuthMiddlewareHandler()
	// if err != nil {
	// 	return fmt.Errorf("error creating auth middleware handler: %v", err)
	// }

	app.Use(logger.New())
	// app.Use("/students", authMiddlewareHandler)

	return nil
}

func registerRoutes(app *fiber.App) error {
	app.Route("/students", func(router fiber.Router) {
		studentHandler := controllers.NewStudentHandler()
		router.Get("/", studentHandler.GetStudents)
		router.Get("/:npm", studentHandler.GetStudentById)
		router.Post("/", studentHandler.PostStudent)
		router.Delete("/:npm", studentHandler.DeleteStudent)
	})

	return nil
}
