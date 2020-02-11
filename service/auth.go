package service

import (
	"errors"
	"fmt"
	"net/http"
	"socket/util"
	"strings"
)

/**
 * 从请求中获取token
 */
func GetTokenData(r *http.Request) (string,error) {
	token := r.Header.Get("Authorization")
	if len(token) <= 0 {
		token = r.FormValue("token")
		if len(token) <= 0 {
			return "",errors.New("token缺失")
		}
	}else{
		tokenSlice := strings.Split(token," ")
		token = tokenSlice[1]
		if len(token) <= 0 {
			return "",errors.New("token缺失")
		}
	}

	data,err := CheckToken(token)
	if err != nil {
		return "",err
	}

	return data,nil
}

/**
 * 校验token
 */
func CheckToken(token string) (string,error) {
	var jwt util.Jwt
	jwt.Token = token
	parseErr := jwt.ParseToken()
	if parseErr != nil {
		util.Error(fmt.Sprintf("解析token失败：%v\n",parseErr.Error()))
		return "",errors.New("非法token")
	}

	if len(jwt.Email.(string)) <= 0 {
		util.Error(fmt.Sprintf("解析token失败：%v\n",token))
		return "",errors.New("非法token201")
	}

	return jwt.Email.(string),nil
}
