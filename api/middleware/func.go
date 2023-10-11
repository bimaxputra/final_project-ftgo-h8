package middleware

import (
	"final_project-ftgo-h8/api/dto"
	"final_project-ftgo-h8/helper"
	"os"

	"github.com/labstack/echo/v4"
)

func (a *authenticationMiddleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// init token from request header
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == ""{
			return dto.WriteResponse(c,401,"no token")
		}

		// init claim
		secretsign := []byte(os.Getenv("SECRETSIGN"))
		claims,err := helper.ParseJWT(tokenString, secretsign)
		if err != nil {
			return dto.WriteResponseWithDetail(c,401,"unauthorized user",err.Error())
		}
		
		// type assertion
		userId := claims["user_id"].(float64)

		// init user from db
		user,err := a.repository.FindUserById(uint(userId))
		if err != nil {
			return dto.WriteResponseWithDetail(c,401,"undefined user",err.Error())
		}

		// set user to context
		c.Set("user",user)
		
		return next(c)
	}
}