package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// 用于签名和验证 JWT 的密钥，与 siwu 项目一致
var jwtSecret = []byte("siwu-web-service.forumSetJwtSecret_S@mpl3ComplexS3cretK3y")

// Claims 自定义的 Claims 结构体，与 siwu 项目一致
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// ValidateToken 验证并解析 Token
func ValidateToken(tokenString string) (*Claims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
