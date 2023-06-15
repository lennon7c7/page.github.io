package util

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	JWTClaims = "Claims"
	JWTUserID = "UID"
	JWTXToken = "X-Token"
)

// JWTCheck 通过请求头验证签名
func JWTCheck(publicKey *rsa.PublicKey, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}

// JWTIssue 生成签名
func JWTIssue(privateKey *rsa.PrivateKey, claims *jwt.StandardClaims, jwtExpires int64) (token string, err error) {
	if claims.ExpiresAt == 0 {
		claims.Id = uuid.NewV4().String()
		now := time.Now()
		claims.IssuedAt = now.Unix()
		claims.ExpiresAt = now.Add(time.Second * time.Duration(jwtExpires)).Unix()
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS384, claims)
	token, err = tokenClaims.SignedString(privateKey)
	return
}
