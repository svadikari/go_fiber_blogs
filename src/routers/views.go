package routers

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/handlers"
	"go_fiber_blogs/src/middleware"
)

func ViewRouters(routers fiber.Router) {
	routers.Route("/", func(router fiber.Router) {
		router.Get("/", middleware.ValidAuthentication, handlers.HomeViewHandler)
		router.Get("/blog", middleware.ValidAuthentication, func(ctx *fiber.Ctx) error {
			return ctx.Render("blog", nil)
		})
		router.Post("/blog", middleware.ValidAuthentication, handlers.CreateBlogHandler)
		router.Get("/blog/:id", middleware.ValidAuthentication, handlers.RenderBlog)
		router.Post("/blog/:id", middleware.ValidAuthentication, handlers.SaveBlog)
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
	routers.Route("/profile", func(router fiber.Router) {
		router.Get("/", middleware.ValidAuthentication, handlers.RenderProfile)
		router.Post("/", middleware.ValidAuthentication, handlers.UpdateProfile)
	})
}
