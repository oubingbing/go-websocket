package util

import (
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	Email interface{}
	secretKey []byte
	Token string
	Exp int64
}

/**
 * 解析token
 */
func (this *Jwt) ParseToken() (error) {
	var configErr error
	this.secretKey,configErr = GetSignKey()
	if configErr != nil {
		return  configErr
	}
	tokenPoint,err := jwt.Parse(this.Token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			Error(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return this.secretKey, nil
	})

	if err != nil{
		return err
	}

	if c, ok := tokenPoint.Claims.(jwt.MapClaims); ok && tokenPoint.Valid {
		this.Email = c["email"]
		return nil
	} else {
		return  err
	}
}

/**
 * 创建token
 */
func (this *Jwt) CreateToken() (error) {
	var configErr error
	this.secretKey,configErr = GetSignKey()
	if configErr != nil {
		return  configErr
	}

	//可以在里面自定义自己需要传输的信息，不要存放机密信息，如密码之类的信息
	type MyCustomClaims struct {
		Email interface{} `json:"email"`//邮箱，用邮箱标记用户信息
		jwt.StandardClaims
	}

	exp := time.Now().Unix()+(3600*24)
	claims := MyCustomClaims{
		this.Email,
		jwt.StandardClaims{
			ExpiresAt:exp,//过期时间，一天
		},
	}

	this.Exp = exp
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)//指定签名方法
	tokenString,err := token.SignedString(this.secretKey)
	if err != nil{
		return  err
	}else{
		this.Token = tokenString
		return nil
	}
}