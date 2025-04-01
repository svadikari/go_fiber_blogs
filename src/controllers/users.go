package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_fiber_blogs/src/database"
	"go_fiber_blogs/src/dtos"
	"go_fiber_blogs/src/middleware"
	"go_fiber_blogs/src/models"
	"go_fiber_blogs/src/utils"
	"strings"
	"time"
)

// CreateUser Create User
//
//	@Summary		Creating User Details
//	@Description	Creating User Details
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.User	true	"user details"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Router			/api/users [post]

func (c *Controller) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(user); errs != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errs, ","),
		}
	}
	db := database.DB.Db
	err := db.Create(&user).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}
	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers GetAllUsers
//
//	@Summary		Getting All Users
//	@Description	Getting All Users
//	@Tags			Users
//	@Produce		json
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Success		200	{object}	[]models.User
//	@Failure		400	{object}	dtos.ErrorResponse
//	@Failure		404	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Router			/api/users [get]
func (c *Controller) GetUsers(ctx *fiber.Ctx) error {
	var users []models.User
	db := database.DB.Db
	db.Find(&users)
	return ctx.Status(fiber.StatusOK).JSON(users)
}

// GetUser Getting User by id
//
//	@Summary		Getting User by id
//	@Description	Getting User by id in detail
//	@Tags			Users
//	@Produce		json
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Param			id	path		string	true	"id of User"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	dtos.ErrorResponse
//	@Failure		404	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Router			/api/users/{id} [get]
func (c *Controller) GetUser(ctx *fiber.Ctx) error {
	userId, _ := ctx.ParamsInt("id")
	var user models.User
	db := database.DB.Db
	db.Find(&user, "id=?", userId)
	if db.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError))
	}
	if user.Id == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}

// UpdateUser Update User by id
//
//	@Summary		Update User by id
//	@Description	Update User by id in detail
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Param			id		path		uint		true	"id of User"
//	@Param			request	body		models.User	true	"Request of Updating User Object"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		404		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Router			/api/users/{id} [put]

func (c *Controller) UpdateUser(ctx *fiber.Ctx) error {
	userId, _ := ctx.ParamsInt("id")
	db := database.DB.Db
	dbUser := new(models.User)
	db.Find(&dbUser, "id=?", userId)
	if dbUser.Id == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	inputUser := new(models.User)
	if err := ctx.BodyParser(inputUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest))
	}
	dbUser.FistName = inputUser.FistName
	dbUser.LastName = inputUser.LastName
	dbUser.Phone = inputUser.Phone
	dbUser.UpdatedAt = time.Now()
	db.Save(&dbUser)
	return ctx.Status(fiber.StatusOK).JSON(dbUser)
}

// DeleteUser Deleting User by id
//
//	@Summary		Deleting User by id
//	@Description	Deleting User by id in detail
//	@Tags			Users
//	@Param			Authorization	header		string		true	"Bearer Token"
//	@Param			id	path	string	true	"id of User"
//	@Success		204
//	@Failure		400	{object}	dtos.ErrorResponse
//	@Failure		404	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Router			/api/users/{id} [delete]
func (c *Controller) DeleteUser(ctx *fiber.Ctx) error {
	userId, _ := ctx.ParamsInt("id")
	user := new(models.User)
	db := database.DB.Db
	err := db.Delete(&user, "id=?", userId).Error
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.NewError(fiber.StatusNotFound))
	}
	return ctx.Status(fiber.StatusNoContent).JSON(fiber.NewError(fiber.StatusNoContent))
}

// GenerateToken Generating JWT Token
//
//	@Summary		Generating JWT Token by username/password
//	@Description	Generating JWT Token by username/password
//	@Tags			Tokens
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dtos.LoginRequest	true	"username/password of User"
//	@Success		200	{object}	dtos.TokenResponse
//	@Failure		400	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Router			/api/generate-token [post]
func (c *Controller) GenerateToken(ctx *fiber.Ctx) error {

	userCredentials := new(dtos.LoginRequest)
	if err := ctx.BodyParser(userCredentials); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	validator := middleware.XValidator{}
	if errs := validator.Validate(userCredentials); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest,
			"errors": errs})
	}
	db := database.DB.Db
	var dbUser models.User
	err := db.Where("user_name = ?", userCredentials.UserName).First(&dbUser).Error
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest,
			"No User Found! Please create an account first!"))
	}
	if err := utils.VerifyPassword(userCredentials.Password, dbUser.Password); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Invalid Credentials!"))
	}
	token, err := middleware.CreateToken(&dbUser)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	tokenResponse := new(dtos.TokenResponse)
	tokenResponse.Token = token
	return ctx.Status(fiber.StatusOK).JSON(tokenResponse)
}
