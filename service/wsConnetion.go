package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"socket/util"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	clientId  string
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex //对closeChan关闭上锁
	isClosed  bool       // 防止closeChan被关闭多次
}

/**
 * 初始化连接
 */
func InitConnect(wsConn *websocket.Conn, clientId string) *Connection {
	conn := &Connection{
		wsConn:    wsConn,
		clientId:  clientId,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
		isClosed:  false,
	}

	go conn.WriteMessageFromChan()

	return conn
}

/**
 * 推送消息到发送chan
 */
func (conn *Connection) PushToChan(data []byte) error {
	var err error
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		//关闭渠道和连接
		conn.Close()
		err = errors.New("连接已关闭")
	}

	return err
}

/**
 * 推动数据到客户端
 */
func (conn *Connection) WriteMessageFromChan() error {
	//从channel不停的去处需要发送的数据
	var data []byte
	var err error
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			//渠道关闭，关闭所有渠道
			conn.Close()
		}
		err = conn.wsConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			//推送数据失败，标记连接关闭
			util.Error(fmt.Sprintf("推送数据失败：%v\n", err.Error()))
			conn.Close()
			break
		}
	}

	return err
}

/**
 * 关闭长连接
 */
func (conn *Connection) Close() {
	conn.wsConn.Close()

	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}
