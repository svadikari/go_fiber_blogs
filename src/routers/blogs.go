package routers

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/controllers"
	"go_fiber_blogs/src/middleware"
)

func BlogApiRouters(router fiber.Router) {
	blogRouter := router.Group("/blogs", middleware.ValidateApiAuth)
	blogRouter.Post("/", controllers.CreateBlog)
	blogRouter.Get("/", controllers.GetBlogs)
	blogRouter.Get("/:id", controllers.GetBlog)
	blogRouter.Put("/:id", controllers.UpdateBlog)
	blogRouter.Delete("/:id", controllers.DeleteBlog)
}
