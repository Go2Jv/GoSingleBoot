package jwt

import (
	"GoSingleBoot/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Issuer string `json:"issuer"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	var claims = Claims{
		UserID: userID,
		Issuer: config.Config.JwtCfg.Issuer,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.Config.JwtCfg.Expiration))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Config.JwtCfg.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("非法Token签名")
		}
		return config.Config.JwtCfg.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, errors.New("非法Token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("解析Token失败")
	}
	return claims, nil
}

func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	if err != nil {
		return false
	}
	return true
}
