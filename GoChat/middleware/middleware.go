package middleware

import (
	"GoChat/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey    = "authorization"
	authorizationHeaderPrefix = "bearer"
	authorizationPayloadKey   = "authorization_payload"
	accessTokenQueryKey       = "token"
)

func extractAccessToken(ctx *gin.Context) (string, error) {
	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		if token := ctx.Query(accessTokenQueryKey); token != "" {
			return token, nil
		}
		return "", errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationHeaderPrefix {
		return "", errors.New("unsupported authorization type")
	}

	return fields[1], nil
}

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := extractAccessToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func GetAuthPayload(ctx *gin.Context) (*token.Payload, bool) {
	value, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, false
	}

	payload, ok := value.(*token.Payload)
	return payload, ok
}
