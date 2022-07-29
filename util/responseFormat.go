package util

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type Response struct {
	ErrorCode    int         `json:"code"`
	ErrorMessage string      `json:"message"`
	Data         interface{} `json:"data"`
}

func (r *Response) ResponseError() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	result, err := json.Marshal(r)
	return result, err
}

func (r *Response) ResponseSuccess() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	result, err := json.Marshal(r)
	return result, err
}

func ResponseJson(w http.ResponseWriter, code int, message string, data interface{}) {
	var (
		res  Response
		json []byte
		err  error
	)
	res.ErrorCode = code
	res.ErrorMessage = message
	res.Data = data

	if code != http.StatusOK {
		json, err = res.ResponseError()
	} else {
		json, err = res.ResponseSuccess()
	}

	if err != nil {
		fmt.Printf("json转换出错：%v\n", err)
	}

	w.Write(json)
}

func Json(code int, message string, data interface{}) []byte {
	var (
		res  Response
		json []byte
		err  error
	)
	res.ErrorCode = code
	res.ErrorMessage = message
	res.Data = data

	json, err = res.ResponseSuccess()

	if err != nil {
		fmt.Printf("json转换出错：%v\n", err)
	}

	return json
}
