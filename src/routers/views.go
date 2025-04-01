package routers

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/handlers"
	"go_fiber_blogs/src/middleware"
)

func ViewRouters(routers fiber.Router) {
	routers.Route("/", func(router fiber.Router) {
		router.Get("/", middleware.ValidAuthentication, handlers.HomeViewHandler)
		router.Get("/create", middleware.ValidAuthentication, func(ctx *fiber.Ctx) error {
			return ctx.Render("new_blog", nil)
		})
		router.Post("/create", middleware.ValidAuthentication, handlers.CreateBlogHandler)
	})
	routers.Route("/register", func(router fiber.Router) {
		router.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Render("register", nil)
		})
		router.Post("/", handlers.Register)
	})
	routers.Route("/login", func(router fiber.Router) {
		router.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Render("login", nil)
		})
		router.Post("/", handlers.Login)
	})
	routers.Route("/logout", func(router fiber.Router) {
		router.Get("/", func(ctx *fiber.Ctx) error {
			ctx.ClearCookie()
			ctx.Locals("userId", 0)
			return ctx.Render("login", nil)
		})
	})
}
