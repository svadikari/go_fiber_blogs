package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go_fiber_blogs/src/database"
	"go_fiber_blogs/src/dtos"
	"go_fiber_blogs/src/middleware"
	"go_fiber_blogs/src/models"
	"go_fiber_blogs/src/utils"
	"time"
)

func HomeViewHandler(ctx *fiber.Ctx) error {
	var blogs []models.Blog
	db := database.DB.Db
	err := db.Model(&models.Blog{}).Order("updated_at DESC").Preload("Author").Find(&blogs).Error
	if err != nil {
		log.Info(err.Error())
	}
	return ctx.Render("home", fiber.Map{"blogs": blogs})
}

func CreateBlogHandler(ctx *fiber.Ctx) error {
	newBlog := new(dtos.BlogRequest)
	if err := ctx.BodyParser(newBlog); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(newBlog); errs != nil {
		return ctx.Render("new_blog", fiber.Map{"errors": errs})
	}
	blog := models.Blog{}
	blog.Title = newBlog.Title
	blog.Content = newBlog.Content
	blog.AuthorId = ctx.Locals("userId").(uint64)
	db := database.DB.Db
	err := db.Create(&blog).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return ctx.Redirect("/", 302)
}

func Register(ctx *fiber.Ctx) error {
	type RegistrationForm struct {
		dtos.UserRequest
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
		return ctx.Render("login", fiber.Map{"errors": err})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:  "jwtToken",
		Value: token,
	})
	ctx.Locals("userId", dbUser.Id)
	return ctx.Redirect("/", fiber.StatusFound)
}

func RenderProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(uint64)
	db := database.DB.Db
	user := models.User{}
	if err := db.Find(&user, "id=?", userId).Error; err != nil {
		log.Info(err.Error())
	}
	return ctx.Render("profile", fiber.Map{"FirstName": user.FistName, "LastName": user.LastName, "Phone": user.Phone})
}

func UpdateProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(uint64)
	updateUser := new(dtos.UserProfile)

	if err := ctx.BodyParser(updateUser); err != nil {
		return ctx.Render("profile", fiber.Map{"errors": err.Error(),
			"FirstName": updateUser.FistName, "LastName": updateUser.LastName, "Phone": updateUser.Phone})
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(updateUser); errs != nil {
		return ctx.Render("profile", fiber.Map{"errors": errs,
			"FirstName": updateUser.FistName, "LastName": updateUser.LastName, "Phone": updateUser.Phone})
	}

	db := database.DB.Db
	var dbUser models.User
	_ = db.Where("id = ?", userId).First(&dbUser).Error
	if dbUser.Id == 0 {
		return ctx.Render("profile", fiber.Map{"errors": []string{"User does not exist, please choose another one"}})
	}

	dbUser.FistName = updateUser.FistName
	dbUser.LastName = updateUser.LastName
	dbUser.Phone = updateUser.Phone
	dbUser.UpdatedAt = time.Now()

	err := db.Updates(&dbUser).Error
	if err != nil {
		return ctx.Render("profile", fiber.Map{"errors": []string{err.Error()}})
	}
	return ctx.Redirect("/", fiber.StatusFound)
}

func RenderBlog(ctx *fiber.Ctx) error {
	blogId, _ := ctx.ParamsInt("id")
	var dbBlog models.Blog
	db := database.DB.Db
	_ = db.Where("id = ?", blogId).First(&dbBlog).Error
	if dbBlog.Id == 0 {
		return ctx.Redirect("/", fiber.StatusFound)
	}
	return ctx.Render("blog", fiber.Map{"Id": dbBlog.Id, "Title": dbBlog.Title, "Content": dbBlog.Content})
}

func SaveBlog(ctx *fiber.Ctx) error {
	blogId, _ := ctx.ParamsInt("id")
	blogRequest := new(dtos.BlogRequest)

	if err := ctx.BodyParser(blogRequest); err != nil {
		return ctx.Render("blog", fiber.Map{"errors": err.Error(),
			"Id": blogId, "Title": blogRequest.Title, "Content": blogRequest.Content})
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(blogRequest); errs != nil {
		return ctx.Render("blog", fiber.Map{"errors": errs,
			"Id": blogId, "Title": blogRequest.Title, "Content": blogRequest.Content})
	}

	db := database.DB.Db
	var dbBlog models.Blog
	_ = db.Where("id = ?", blogId).First(&dbBlog).Error
	if dbBlog.Id == 0 {
		return ctx.Render("blog", fiber.Map{"errors": []string{"Blog does not exist, please choose another one"},
			"Id": blogId, "Title": blogRequest.Title, "Content": blogRequest.Content})
	}

	dbBlog.Title = blogRequest.Title
	dbBlog.Content = blogRequest.Content
	dbBlog.UpdatedAt = time.Now()

	err := db.Updates(&dbBlog).Error
	if err != nil {
		return ctx.Render("blog", fiber.Map{"errors": []string{err.Error()}, "Id": blogId, "Title": blogRequest.Title,
			"Content": blogRequest.Content})
	}
	return ctx.Redirect("/", fiber.StatusFound)
}
