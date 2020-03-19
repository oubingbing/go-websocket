## 部署流程

### 一、安装golang环境，需要1.13以上的版本
	wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
	tar -C /usr/local/ -xvf go1.14.linux-amd64.tar.gz
	vim /etc/profile
	export GOROOT=/usr/local/go
	export PATH=$PATH:/usr/local/go/bin
	export GOPATH=/data/golang
	
	go version

### 二、拉取代码
发布地址：ssh://git@git.galaxymx.com:22333/zypt/websocket.git (git pull)
发布分支： master
发布备注：

### 三、设置mod环境变量
	vim /etc/profile
	GOPROXY=https://goproxy.io
	GO111MODULE=on

### 四、打包
在项目目录下执行一下命令

	set GOOS=linux
	set GOARCH=amd64
	go build

### 五、启动服务
先给生成的二级制文件scoket读写的权限
然后执行以下命令启动服务

`nohup ./socket &`

该服务的端口为 8088

### 六、配置nginx
参考配置

	[root@txy-op-elk-05 conf.d]# vim socket.galaxymx.com.conf
	server {
		listen  80;
		server_name   socket.galaxymx.com;
	
		charset    utf-8;
		default_type 'text/html';
		client_max_body_size 4G;
	
		error_log /data/logs/nginx/socket.galaxymx.com.error.log;
		access_log /data/logs/nginx/socket.galaxymx.com.access.log;
		index index.php index.html;
	
		location / {
			proxy_pass http://127.0.0.1:8088;
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection "upgrade";
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		}
	}
	
完成部署

### api文档

#### 1.授权获得访问长连接权限

**请求URL：**
- ` http://127.0.0.1:8088/auth`

**请求方式：**
- POST

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|sign |是  |string |后端返回的token,aes-128-ecb对称加密字符串|


 **返回示例**

```
{
    "code":200,
    "message":"授权成功",		"data":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IloxOHpYMkZmTVE9PSIsImV4cCI6MTU4MTQ5NTgyOX0.fUUvADIm7otda4Z58fq1PIgcXnzBGhx6_ERDBBQjYKc",
    "contact_email":"875307054@qq.com"
}
```

 **返回参数说明** 

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|code |  int |不为200即报错 |
|message |  string |接口提示信息 |
|data |  string |token |

#### 1.握手建立长连接

**请求URL：**
- ` ws://127.0.0.1:8088/ws`

**请求方式：**
- GET

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|token |是  |string |后端返回的token|


 **返回示例**

```
{"code":201,"message":"连接成功","data":null,"contact_email":"875307054@qq.com"}

```

 **返回参数说明** 

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|code |  int |201为心跳检测或者握手成功，200为业务数据推送 |
|message |  string |接口提示信息 |
|data |  string |业务数据 |



#### 2.握手建立长连接

**请求URL：**
- ` ws://127.0.0.1:8088/ws`

**请求方式：**
- GET

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|token |是  |string |后端返回的token|


 **返回示例**

```
{"code":201,"message":"连接成功","data":null,"contact_email":"875307054@qq.com"}
```

 **返回参数说明** 

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|code |  int |201为心跳检测或者握手成功，200为业务数据推送 |
|message |  string |接口提示信息 |
|data |  string |业务数据 |


#### 3.消息推送

**请求URL：**
- ` http://127.0.0.1:8088/push`

**请求方式：**
- POST

**参数：**

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
|message |是  |string |推送的业务消息|
|token |是  |string |bearer token，或者可以跟在URL地址后，如：http://127.0.0.1:8088/push?token=******|


 **返回示例**

```
{
  "code" : 200
  "message" : "推送成功"
  "data" => null
  "contact_email" : "875307054@qq.com"
}

```

 **返回参数说明** 

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|code |  int |不为200即报错 |
|message |  string |接口提示信息 |
|data |  string |nil |
