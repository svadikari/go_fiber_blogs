package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go_fiber_blogs/src/database"
	"go_fiber_blogs/src/middleware"
	"go_fiber_blogs/src/models"
	"go_fiber_blogs/src/utils"
)

func HomeViewHandler(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Locals("userId"))
	var blogs []models.Blog
	db := database.DB.Db
	err := db.Find(&blogs).Error
	if err != nil {
		log.Info(err.Error())
	}
	return ctx.Render("home", fiber.Map{"Title": "Home", "Blogs": blogs})
}

func CreateBlogHandler(ctx *fiber.Ctx) error {
	blog := new(models.Blog)
	if err := ctx.BodyParser(blog); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(blog); errs != nil {
		return ctx.Render("new_blog", fiber.Map{"errors": errs})
	}
	db := database.DB.Db
	err := db.Create(&blog).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return ctx.Redirect("/", 302)
}

func Register(ctx *fiber.Ctx) error {
	type RegistrationForm struct {
		models.User
		ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=32"`
	}
	registerForm := new(RegistrationForm)

	if err := ctx.BodyParser(registerForm); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(registerForm); errs != nil {
		return ctx.Render("register", fiber.Map{"errors": errs})
	}

	db := database.DB.Db
	var dbUser models.User
	_ = db.Where("user_name = ?", registerForm.UserName).First(&dbUser).Error
	if dbUser.Id != 0 {
		return ctx.Render("register", fiber.Map{"errors": []string{"User already exists, please choose another one"}})
	}

	var user models.User
	user.UserName = registerForm.UserName

	user.Password = utils.GenerateHash(registerForm.Password)
	user.FistName = registerForm.FistName
	user.LastName = registerForm.LastName
	user.Phone = registerForm.Phone

	err := db.Create(&user).Error

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return ctx.Redirect("/login", fiber.StatusTemporaryRedirect)
}
func Login(ctx *fiber.Ctx) error {
	type LoginForm struct {
		UserName string `json:"userName" form:"username" validate:"required,min=5,max=10" gorm:"user_name"`
		Password string `json:"password" form:"password" validate:"required,min=5,max=10" gorm:"password"`
	}
	loginForm := new(LoginForm)
	if err := ctx.BodyParser(loginForm); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(loginForm); errs != nil {
		return ctx.Render("login", fiber.Map{"errors": errs})
	}
	db := database.DB.Db
	var dbUser models.User
	err := db.Where("user_name = ?", loginForm.UserName).First(&dbUser).Error
	if err != nil {
		return ctx.Render("login", fiber.Map{"error": "user not found"})
	}
	if err := utils.VerifyPassword(loginForm.Password, dbUser.Password); err != nil {
		return ctx.Render("login", fiber.Map{"error": "invalid password"})
	}
	token, err := middleware.CreateToken(&dbUser)
	if err != nil {
		fmt.Println(err)
		return ctx.Render("login", fiber.Map{"errors": err})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:  "jwtToken",
		Value: token,
	})
	ctx.Locals("userId", dbUser.Id)
	return ctx.Redirect("/", fiber.StatusFound)
}
