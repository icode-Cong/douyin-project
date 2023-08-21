package tokenservice

import (
	"context"
	"tokenService/protos/token"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
}

var jwtSecret = []byte("1122233")

type Claims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func (*TokenService) ParseTokenToId(ctx context.Context, request *token.ParseTokenToIdRequest, response *token.ParseTokenToIdResponse) error {
	token := request.Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			response.UserId = int32(claims.Id)
			return nil
		}
	}
	return err
}
