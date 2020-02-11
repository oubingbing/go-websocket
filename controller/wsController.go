package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"socket/service"
	"socket/util"
)

type Ws struct {
	GlobalSocket *service.GlobalSocket
}

/**
 * 授权，需要进行签名校验
 */
func (ws *Ws) Auth(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		util.ResponseJson(w,502,"method error",nil)
		return
	}

	encryptString := r.PostFormValue("sign")
	key := []byte("7yY2tYZdPuNSBVU9")
	if len(encryptString) <= 0 {
		util.ResponseJson(w,503,"签名错误",nil)
		return
	}

	base64Decode,base64Err:= base64.StdEncoding.DecodeString(encryptString)
	if base64Err != nil {
		util.ResponseJson(w,501,"签名错误",nil)
		return
	}

	decryptCode,decodeErr := util.AesDecryptECB(base64Decode, key)
	if decodeErr != nil {
		util.ResponseJson(w,502,"签名错误",nil)
		return
	}

	//发放token
	var jwt util.Jwt
	jwt.Email = decryptCode
	err := jwt.CreateToken()
	if err != nil{
		util.Error(fmt.Sprintf("创建token失败：%v\n",err.Error()))
		util.ResponseJson(w,504,"签名错误",nil)
		return
	}

	util.ResponseJson(w,200,"授权成功",jwt.Token)
}

/**
 * 处理协议升级
 */
func (ws *Ws) UpGrad(w http.ResponseWriter, r *http.Request)  {
	clientId,err := service.GetTokenData(r)
	if err != nil {
		util.ResponseJson(w,505,err.Error(),nil)
		fmt.Println(err.Error())
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn,err := upgrader.Upgrade(w,r,nil)
	if err != nil{
		//连接错误
		fmt.Printf("升级协议失败：%v\n",err)
		return
	}

	var clientConn *service.Connection
	clientConn = service.InitConnect(wsConn,clientId)
	ws.GlobalSocket.ClientConnMap[clientId] = clientConn

	jsonData := util.Json(http.StatusOK,"连接成功",nil)
	err = ws.GlobalSocket.ClientConnMap[clientId].PushToChan(jsonData)
	if err != nil {
		fmt.Printf("推送消息错误：%v\n",err.Error())
	}
}

/**
 * 推送消息到指定客户端
 */
func (ws *Ws) Push(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		util.ResponseJson(w,502,"method error",nil)
		return
	}

	data := r.PostFormValue("message")
	clientId,err := service.GetTokenData(r)
	if err != nil {
		util.ResponseJson(w,506,"无权访问",nil)
		return
	}

	socketConn,ok := ws.GlobalSocket.ClientConnMap[clientId]
	if !ok {
		util.ResponseJson(w,507,"客户端不存在",nil)
		return
	}

	jsonData := util.Json(http.StatusOK,data,nil)
	err = socketConn.PushToChan(jsonData)
	if err != nil {
		util.ResponseJson(w,508,"推送失败",nil)
		return
	}

	util.ResponseJson(w,200,"推送成功",nil)
	return
}