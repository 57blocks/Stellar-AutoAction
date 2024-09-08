package middleware

import (
	"errors"
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	"github.com/57blocks/auto-action/server/internal/pkg/logx"
	"github.com/57blocks/auto-action/server/internal/service/cs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SecretKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing api key"})
			return
		}

		apiKey, err := cs.Conductor.APIKey(c)
		if err != nil {
			c.Error(err)
			return
		}

		if apiKey != key {
			c.Error(errorx.UnauthorizedWithMsg("invalid api key"))
			return
		}

		c.Next()
	}
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			return
		}

		jwtToken, err := jwtx.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		claimMap, ok := jwtToken.Claims.(jwt.MapClaims) // TODO: remove the type conversion after testing
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		// TODO: make the keys into constant type
		// TODO: confirm that if all standard claims should be included
		c.Set("jwt_raw", jwtToken.Raw)
		//c.Set("jwt_sub", claimMap["sub"])
		//c.Set("jwt_iat", claimMap["iat"])
		//c.Set("jwt_iss", claimMap["iss"])
		c.Set("jwt_exp", claimMap["exp"])
		c.Set("jwt_account", claimMap["account"])
		c.Set("jwt_organization", claimMap["organization"])
		c.Set("jwt_environment", claimMap["environment"])

		logx.Logger.DEBUG("authentication success")

		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		logx.Logger.DEBUG("authorization success")

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Here are the post-middlewares
// Their execution order is the opposite of the register order.

func PostResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		c.Header("Content-Type", "application/json")
	}
}

func PostHandleErr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		e := c.Errors.Last()

		err := new(errorx.Errorx)
		if errors.As(e.Err, &err) {
			c.JSON(err.Status(), &struct{ Error interface{} }{Error: err})
			return
		}

		// unrecognized error
		c.JSON(http.StatusInternalServerError, &struct{ Error interface{} }{Error: e})
		return
	}
}
