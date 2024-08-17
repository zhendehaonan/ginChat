package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("jilinjianke")

type Claims struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

// 发放token
func GenerateToken(id uint, userName string, authority int) (string, error) {
	//发放时间
	now := time.Now()
	//过期时间24小时
	add := now.Add(24 * time.Hour)
	claims := Claims{
		ID:        id,
		UserName:  userName,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: add.Unix(),
			Issuer:    "Jian-Ming",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	fmt.Println("jwt: ", jwtSecret)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

/*
// 邮箱token
type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

// 签发邮箱token
func GenerateEmailToken(userId, Operation uint, email, password string) (string, error) {
	//发放时间
	now := time.Now()
	//过期时间24小时
	add := now.Add(24 * time.Hour)
	emailClaims := EmailClaims{
		UserID:        userId,
		Email:         email,
		Password:      password,
		OperationType: Operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: add.Unix(),
			Issuer:    "FanOne_Mall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, emailClaims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 解析邮箱token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	fmt.Println("jwt: ", jwtSecret)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}*/
