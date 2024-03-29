package middleware

import (
	"authorization_service/internal/pkg/jwt-service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(ctx *gin.Context) {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		var ok bool
		tokenHeader, ok = ctx.GetQuery("token")
		if !ok || tokenHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is not provided"})
			return
		}
	}
	header := strings.Split(tokenHeader, " ")
	if len(header) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token is invalid format"})
		return
	}

	access := header[1]
	claims, ok := jwtService.VerifyToken(access)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token is invalid format"})
		return
	}
	ctx.Set("userClaims", claims)
	ctx.Next()
}
