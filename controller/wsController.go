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
		util.Error(fmt.Sprintf("ethod error-5002"))
		util.ResponseJson(w,500,"method error",nil)
		return
	}

	encryptString,ok:= util.Input(r,"sign")
	key := []byte("7yY2tYZdPuNSBVU9")
	if !ok {
		util.Error(fmt.Sprintf("签名错误-5003"))
		util.ResponseJson(w,500,"签名错误",nil)
		return
	}

	base64Decode,base64Err:= base64.StdEncoding.DecodeString(encryptString.(string))
	if base64Err != nil {
		util.Error(fmt.Sprintf("签名错误-5001"))
		util.ResponseJson(w,500,"签名错误",nil)
		return
	}

	decryptCode,decodeErr := util.AesDecryptECB(base64Decode, key)
	if decodeErr != nil {
		util.Error(fmt.Sprintf("签名错误-5002"))
		util.ResponseJson(w,500,"签名错误",nil)
		return
	}

	//发放token
	var jwt util.Jwt
	jwt.Email = decryptCode
	err := jwt.CreateToken()
	if err != nil{
		util.Error(fmt.Sprintf("创建token失败-5004：%v\n",err.Error()))
		util.ResponseJson(w,500,"签名错误",nil)
		return
	}

	mp := make(map[string]interface{})
	mp["token"] = jwt.Token
	mp["expire"] = jwt.Exp
	util.ResponseJson(w,200,"授权成功",mp)
}

/**
 * 处理协议升级
 */
func (ws *Ws) UpGrad(w http.ResponseWriter, r *http.Request)  {
	clientId,err := service.GetTokenData(r)
	if err != nil {
		util.Error(fmt.Sprintf("获取token信息失败-5005：%v\n",err.Error()))
		util.ResponseJson(w,500,err.Error(),nil)
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
		util.Error(fmt.Sprintf("升级协议失败：%v\n",err.Error()))
		return
	}

	var clientConn *service.Connection
	clientConn = service.InitConnect(wsConn,clientId)
	ws.GlobalSocket.ClientConnMap[clientId] = clientConn

	jsonData := util.Json(201,"连接成功",nil)
	err = ws.GlobalSocket.ClientConnMap[clientId].PushToChan(jsonData)
	if err != nil {
		util.Error(fmt.Sprintf("推送消息错误：%v\n",err.Error()))
	}
}

/**
 * 推送消息到指定客户端
 */
func (ws *Ws) Push(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		util.ResponseJson(w,500,"method error",nil)
		return
	}

	data,ok := util.Input(r,"message")
	if !ok {
		util.ResponseJson(w,500,"message不能为空",nil)
		return
	}

	clientId,err := service.GetTokenData(r)
	if err != nil {
		util.ResponseJson(w,500,"无权访问",nil)
		return
	}

	socketConn,ok := ws.GlobalSocket.ClientConnMap[clientId]
	if !ok {
		util.Error(fmt.Sprintf("客户端不存在：%v\n",clientId))
		util.ResponseJson(w,500,"客户端不存在",nil)
		return
	}

	jsonData := util.Json(http.StatusOK,"服务端消息推送",data)
	err = socketConn.PushToChan(jsonData)
	if err != nil {
		util.Error(fmt.Sprintf("推送消息错误：%v\n",err.Error()))
		util.ResponseJson(w,500,"推送失败",nil)
		return
	}

	util.ResponseJson(w,200,"推送成功",nil)
	return
}