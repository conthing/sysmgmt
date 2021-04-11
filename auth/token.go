package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var verifyKey = []byte("conthing")

// ExpireTime 过期时间会换算成秒
const ExpireTime = 3600

type userClaims struct {
	Username    string `json:"username"`
	TokenRandom int    `json:"tokenrandom"`
	jwt.StandardClaims
}

// GenerateToken 通过 UserID 来签发
func GenerateToken(userName string, tokenRandom int) (string, error) {
	// Create the Claims
	claims := userClaims{
		userName,
		tokenRandom,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(verifyKey)
	return tokenString, err
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*userClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("ParseToken failed: %w", err)
	}

	claims := token.Claims.(*userClaims)
	return claims, token.Claims.Valid()
}
