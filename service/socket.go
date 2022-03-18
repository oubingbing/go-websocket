package service

import (
	"fmt"
	"socket/util"
	"time"
)

type GlobalSocket struct {
	ClientConnMap map[string][]*Connection
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
func (socket *GlobalSocket) Heartbeat() {
	for {
		for i, conn := range socket.ClientConnMap {

			newArr := []*Connection{}
			jsonData := util.Json(201, "heartbeat", nil)
			for _, item := range conn {
				if err := item.PushToChan(jsonData); err == nil {
					//没有回应的客户端会被清除
					newArr = append(newArr, item)
				}
			}
			if len(newArr) <= 0 {
				//把所有连接删除
				delete(socket.ClientConnMap, i)
			} else {
				socket.ClientConnMap[i] = newArr
			}
		}

		time.Sleep(10 * time.Second)

		//检测客户端连接数
		util.Info(fmt.Sprintf("客户端连接数：%v\n", socket.CountClient()))
	}
}
