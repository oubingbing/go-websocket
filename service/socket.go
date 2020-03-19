package service

import (
	"fmt"
	"socket/util"
	"time"
)

type GlobalSocket struct {
	ClientConnMap map[string]*Connection
}

/**
 * 客户端连接数
 */
func (socket *GlobalSocket) CountClient() int {
	return len(socket.ClientConnMap)
}

/*
 * 检测心跳
 */
func (socket *GlobalSocket) Heartbeat()  {
	for{
		for _,conn := range socket.ClientConnMap {

			jsonData := util.Json(201,"heartbeat",nil)
			if err := conn.PushToChan(jsonData);err != nil{
				//客户端断开连接
				delete(socket.ClientConnMap,conn.clientId)
			}
		}

		time.Sleep(30*time.Second)

		//检测客户端连接数
		fmt.Printf("客户端连接数：%v\n",socket.CountClient())
	}
}