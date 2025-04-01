package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go_fiber_blogs/src/models"
	"go_fiber_blogs/src/utils"
	"strings"
	"time"
)

func CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user": user.UserName,
			"exp":  time.Now().Add(time.Minute * 2).Unix(),
		})
	tokenString, err := token.SignedString([]byte(utils.GetEnvConfig("SECRET_KEY", "secured-secret-key")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetEnvConfig("SECRET_KEY", "secured-secret-key")), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func ValidateApiAuth(ctx *fiber.Ctx) error {
	jwtToken := ctx.Get("Authorization")
	if jwtToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.NewError(fiber.StatusUnauthorized, "Provide valid token"))
	}
	strArr := strings.Split(jwtToken, " ")
	if len(strArr) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.NewError(fiber.StatusUnauthorized, "Provide valid token"))
	}
	if err := VerifyToken(strArr[1]); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.NewError(fiber.StatusUnauthorized, "invalid token"))
	}
	return ctx.Next()

}

func ValidAuthentication(ctx *fiber.Ctx) error {
	jwtToken := ctx.Cookies("jwtToken")
	if jwtToken == "" {
		return ctx.Redirect("/login", fiber.StatusUnauthorized)
	}
	if err := VerifyToken(jwtToken); err != nil {
		return ctx.Redirect("/login", fiber.StatusUnauthorized)
	}
	return ctx.Next()
}
