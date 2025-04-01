package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"go_fiber_blogs/src/database"
	_ "go_fiber_blogs/src/docs"
	"go_fiber_blogs/src/middleware"
	"go_fiber_blogs/src/routers"
	"go_fiber_blogs/src/utils"
)

//	@title			Blogs API
//	@version		1.0
//	@description	This app deals with blog CRUD APIs
//	@termsOfService	https://www.shyam.com/terms/
//	@contact.name	Shyam
//	@contact.email	shyam@shyam.com
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@BasePath		/

func main() {
	database.InitDB()
	app := InitApp()
	PORT := utils.GetEnvConfig("APPLICATION_PORT", "8080")

	log.Fatal(app.Listen(":" + PORT))
}

func InitApp() *fiber.App {

	engine := html.New("./templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		// Pass in Views Template Engine
		Views: engine,
		// Default global path to search for templates (can be overridden when calling Render())
		ViewsLayout: "layouts/base",
		// Global Error Handler
		ErrorHandler: middleware.ErrorHandler,

		// Enables/Disables access to `ctx.Locals()` entries in rendered templates
		// (defaults to false)
		PassLocalsToViews: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080,http://127.0.0.1:8080,https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	viewRouter := app.Group("/")
	routers.ViewRouters(viewRouter)
	apiRouter := app.Group("/api")
	routers.BlogApiRouters(apiRouter)
	routers.UserApiRouters(apiRouter)
	app.Get("/swagger/*", swagger.HandlerDefault)
	return app
}
