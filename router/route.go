package router

import (
	"net/http"
	"socket/controller"
	"socket/service"
)

func Init()  {
	globalSocket := &service.GlobalSocket{make(map[string][]*service.Connection)}
	go globalSocket.Heartbeat()

	ws := &controller.Ws{GlobalSocket:globalSocket}

	http.HandleFunc("/auth",ws.Auth)
	//升级协议
	http.HandleFunc("/ws", ws.UpGrad)
	//推送消息
	http.HandleFunc("/push", ws.Push)
}