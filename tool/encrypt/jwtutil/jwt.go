package jwtutil

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type Claims struct {
	Data map[string]string
	jwt.StandardClaims
}

const (
	ISSUER = "GATEWAY_SERVER"
	KEY    = "67617465776179736572766572d41d8cd98f00b204e9800998ecf8427e"
)

func GetToken(data map[string]string) (string, error) {
	claims := &Claims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			Issuer:    ISSUER,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstr, err := token.SignedString([]byte(KEY))
	return tokenstr, err
}

func Parse(token string) (map[string]string, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, errors.New("claims不是这个类型")
	}
	data := claims.Data
	data["expireAt"] = strconv.FormatInt(claims.ExpiresAt, 10)
	return claims.Data, nil
}
