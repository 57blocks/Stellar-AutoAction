package middleware

import (
	"errors"
	"github.com/57blocks/auto-action/server/internal/constant"
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/service/cs"
	"github.com/57blocks/auto-action/server/internal/third-party/eslint"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/gin-gonic/gin"
)

func SecretKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(constant.APIKey.Str())
		if key == "" {
			c.Error(errorx.UnauthorizedWithMsg("missing api key"))
			c.Abort()
			return
		}

		apiKey, err := cs.ServiceImpl.APIKey(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if apiKey != key {
			c.Error(errorx.UnauthorizedWithMsg("invalid api key"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader(constant.AuthHeader.Str())
		if jwtToken == "" {
			c.Error(errorx.UnauthorizedWithMsg("missing token"))
			c.Abort()
			return
		}

		c.Set(constant.ClaimRaw.Str(), jwtToken)

		c.Next()
	}
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Error(errorx.UnauthorizedWithMsg("missing token"))
			c.Abort()
			return
		}

		jwtClaims, err := jwtx.RS256.Parse(token)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		claimMap, ok := jwtClaims.(*jwtx.AAClaims) // TODO: remove the type conversion after testing
		if !ok {
			c.Error(errorx.UnauthorizedWithMsg("invalid token"))
			c.Abort()
			return
		}

		c.Set(constant.ClaimSub.Str(), claimMap.StdJWTClaims.Subject)
		c.Set(constant.ClaimIss.Str(), claimMap.StdJWTClaims.Issuer)

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
			c.Abort()
			return
		}

		e := c.Errors.Last()

		err := new(errorx.Errorx)
		if errors.Is(e.Err, err) {
			c.JSON(err.Status(), &struct{ Error interface{} }{Error: err})
			c.Abort()
			return
		}

		// unrecognized error
		c.JSON(http.StatusInternalServerError, &struct{ Error interface{} }{Error: e})
		return
	}
}

func RegisterESLintCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.Error(errorx.BadRequest(err.Error()))
			c.Abort()
			return
		}

		for _, fileHeader := range c.Request.MultipartForm.File {
			for _, header := range fileHeader {
				if err := eslint.Check(header); err != nil {
					c.Error(errorx.BadRequest(err.Error()))
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
