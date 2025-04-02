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
			"userId":    user.Id,
			"firstName": user.FistName,
			"exp":       time.Now().Add(time.Minute * 30).Unix(),
		})
	tokenString, err := token.SignedString([]byte(utils.GetEnvConfig("SECRET_KEY", "secured-secret-key")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type Claims struct {
	UserId    uint64 `json:"userId"`
	FirstName string `json:"firstName"`
	jwt.RegisteredClaims
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utils.GetEnvConfig("SECRET_KEY", "secured-secret-key")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return claims, nil
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
	claims, err := VerifyToken(strArr[1])
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.NewError(fiber.StatusUnauthorized, "invalid token"))
	}
	ctx.Locals("userId", claims.UserId)
	ctx.Locals("firstName", claims.FirstName)
	return ctx.Next()

}

func ValidAuthentication(ctx *fiber.Ctx) error {
	jwtToken := ctx.Cookies("jwtToken")
	if jwtToken == "" {
		return ctx.Redirect("/login", fiber.StatusTemporaryRedirect)
	}

	claims, err := VerifyToken(jwtToken)
	if err != nil {
		return ctx.Redirect("/login", fiber.StatusTemporaryRedirect)
	}
	ctx.Locals("userId", claims.UserId)
	ctx.Locals("firstName", claims.FirstName)
	return ctx.Next()
}
