package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/database"
	"go_fiber_blogs/src/dtos"
	"go_fiber_blogs/src/middleware"
	"go_fiber_blogs/src/models"
	"strings"
	"time"
)

// CreateBlog Create Blog
//
//	@Summary		Creating Blog Details
//	@Description	Creating Blog Details
//	@Tags			Blogs
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Param			blog			body		dtos.BlogRequest	true	"blog details"
//	@Success		201				{object}	models.Blog
//	@Failure		400				{object}	dtos.ErrorResponse
//	@Failure		500				{object}	dtos.ErrorResponse
//	@Router			/api/blogs/ [post]

func CreateBlog(ctx *fiber.Ctx) error {
	newBlog := new(dtos.BlogRequest)
	if err := ctx.BodyParser(newBlog); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(newBlog); errs != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errs, ","),
		}
	}
	blog := models.Blog{}
	blog.Title = newBlog.Title
	blog.Content = newBlog.Content
	blog.AuthorId = ctx.Locals("userId").(uint64)
	db := database.DB.Db
	err := db.Create(&blog).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}
	return ctx.Status(fiber.StatusCreated).JSON(blog)
}

// GetBlogs Get All Blogs
//
//	@Summary		Getting All Blogs
//	@Description	Getting All Blogs
//	@Tags			Blogs
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Success		200				{object}	[]models.Blog
//	@Router			/api/blogs/ [get]

func GetBlogs(ctx *fiber.Ctx) error {
	var Blogs []models.Blog
	db := database.DB.Db
	db.Find(&Blogs).Preload("Author").Find(&Blogs)
	return ctx.Status(fiber.StatusOK).JSON(Blogs)
}

// GetBlog Get Blog
//
//	@Summary		Getting a Blog
//	@Description	Getting a Blog
//	@Tags			Blogs
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Param			id				path		int		true	"blog id"
//	@Success		200				{object}	models.Blog
//	@Failure		404				{object}	dtos.ErrorResponse
//	@Failure		500				{object}	dtos.ErrorResponse
//	@Router			/api/blogs/{id} [get]

func GetBlog(ctx *fiber.Ctx) error {
	blogId, _ := ctx.ParamsInt("id")
	var blog models.Blog
	db := database.DB.Db
	db.Find(&blog, "id=?", blogId)
	if db.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError))
	}
	if blog.Id == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	return ctx.Status(fiber.StatusOK).JSON(blog)
}

// UpdateBlog Update Blog
//
//	@Summary		Updating a Blog
//	@Description	Updating a Blog
//	@Tags			Blogs
//	@Produce		json
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Param			id				path		uint		true	"blog id"
//	@Param			blog			body		dtos.BlogRequest	true	"blog details"
//	@Success		200				{object}	models.Blog
//	@Failure		400				{object}	dtos.ErrorResponse
//	@Failure		404				{object}	dtos.ErrorResponse
//	@Failure		500				{object}	dtos.ErrorResponse
//	@Router			/api/blogs/{id} [put]

func UpdateBlog(ctx *fiber.Ctx) error {
	blogId, _ := ctx.ParamsInt("id")
	db := database.DB.Db
	dbBlog := new(models.Blog)
	db.Find(&dbBlog, "id=?", blogId)
	if dbBlog.Id == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	inputBlog := new(dtos.BlogRequest)
	if err := ctx.BodyParser(inputBlog); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	dbBlog.Title = inputBlog.Title
	dbBlog.Content = inputBlog.Content
	dbBlog.UpdatedAt = time.Now()
	db.Save(&dbBlog)
	return ctx.Status(fiber.StatusOK).JSON(dbBlog)
}

// DeleteBlog Delete Blog
//
//	@Summary		Deleting a Blog
//	@Description	Deleting a Blog
//	@Tags			Blogs
//	@Param			Authorization	header	string	true	"Bearer Token"
//	@Param			id				path	uint	true	"blog id"
//	@Success		204
//	@Failure		404	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Router			/api/blogs/{id} [delete]

func DeleteBlog(ctx *fiber.Ctx) error {
	blogId, _ := ctx.ParamsInt("id")
	blog := new(models.Blog)
	db := database.DB.Db
	err := db.Delete(&blog, "id=?", blogId).Error
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	return ctx.Status(fiber.StatusNoContent).JSON(fiber.NewError(fiber.StatusNoContent))
}
