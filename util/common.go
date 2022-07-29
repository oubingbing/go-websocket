package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Input(r *http.Request, key string) (interface{}, bool) {
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	value, ok := data[key]
	fmt.Println(data)
	return value, ok
}
