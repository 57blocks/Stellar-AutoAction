package middleware

import (
	"net/http"

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

func Response() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func Error() gin.HandlerFunc {
	// TODO: finish the general error handler in middleware
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.Header("Content-Type", "application/json")
			lastError := c.Errors.Last()

			status := http.StatusInternalServerError // 默认状态码
			if err, ok := lastError.Err.(interface{ Status() int }); ok {
				status = err.Status()
			}
			c.JSON(
				status,
				gin.H{
					"message": lastError.Error(),
				},
			)

			c.Abort()
		}
	}
}
