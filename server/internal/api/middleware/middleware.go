package middleware

import (
	"errors"
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SecretKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing secret key"})
			return
		}

		secret, err := svcOrg.Conductor.OrgSecret(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if secret != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid secret key"})
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

		pkgLog.Logger.DEBUG("authentication success")

		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		pkgLog.Logger.DEBUG("authorization success")

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
