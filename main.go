package main

import (
	"fmt"
	"net/http"
	"socket/router"
)

func main() {
	router.Init()
	fmt.Println("服务器启动成功")
	http.ListenAndServe("localhost:8088", nil)
}
