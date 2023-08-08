package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32
const iss = "simple_bank"

type JWTMaker struct {
	secretKey string
}

type MyCustomClaims struct {
	Username string
	jwt.RegisteredClaims
}

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("secret key length is short, should be atleast %d", minSecretKeySize)
	}
	return &JWTMaker{secretkey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string,*Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "",payload ,err
	}
	claims := MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			Issuer:    iss,
			ID:        payload.ID.String(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token,err:= jwtToken.SignedString([]byte(maker.secretKey))
	return token,payload,err
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	Keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, Keyfunc)
	if err != nil {
		return nil, err
	}
	if claim, ok := jwtToken.Claims.(*MyCustomClaims); ok && jwtToken.Valid {
		return &Payload{
			ID:        uuid.MustParse(claim.ID),
			Username:  claim.Username,
			IssuedAt:  claim.IssuedAt.Time,
			ExpiredAt: claim.ExpiresAt.Time,
		}, nil
	} else {
		return nil, fmt.Errorf("token can't be parsed")
	}
}
