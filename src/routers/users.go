package routers

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/controllers"
	"go_fiber_blogs/src/middleware"
)

func UserApiRouters(router fiber.Router) {
	c := controllers.NewController()
	router.Post("/generate-token", c.GenerateToken)
	router.Route("/users", func(router fiber.Router) {
		router.Post("/", c.CreateUser)
		router.Get("/", middleware.ValidateApiAuth, c.GetUsers)
		router.Get("/:id", middleware.ValidateApiAuth, c.GetUser)
		router.Put("/:id", middleware.ValidateApiAuth, c.UpdateUser)
		router.Delete("/:id", middleware.ValidateApiAuth, c.DeleteUser)
	})
}
