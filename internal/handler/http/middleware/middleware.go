package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/webping/internal/handler/http/response"
)

const (
	authorizationHeader = "Authorization"
	admin               = "admin"
)

type validationFunc interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

func AdminIdentityMiddleware(service validationFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("empty auth header"))
			c.Abort()
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("invalid with header"))
			c.Abort()
			return
		}

		adminR, err := service.ValidateToken(c, headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err.Error()))
			c.Abort()
			return
		}

		c.Set(admin, adminR)
		c.Next()
	}
}
