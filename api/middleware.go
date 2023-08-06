package api

import (
	"errors"
	"fmt"
	"net/http"
	"simple_bank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header is provided")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s",authorizationType)
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		paylaod,err := tokenMaker.VerifyToken(accessToken)
		fmt.Printf("%+v",paylaod)
		if err != nil{
			err := errors.New("could not verify the token")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey,paylaod)
		ctx.Next()
	}
}
