/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 16:04
 */

package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gowebsocket/helper"
	"gowebsocket/models"
	"net/http"
	"time"
)

var (
	clientManager = NewClientManager() // 管理者
	appIds        = []uint32{101, 102} // 全部的平台

	serverIp   string
	serverPort string
)

func GetAppIds() []uint32 {

	return appIds
}

func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)

	return
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}

	return
}

func InAppIds(appId uint32) (inAppId bool) {

	for _, value := range appIds {
		if value == appId {
			inAppId = true

			return
		}
	}

	return
}

// 启动程序
func StartWebSocket() {

	serverIp = helper.GetServerIp()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort

	// 请求升级为 WebSocket
	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()

	fmt.Println(
		"-------------------------------------------\n"+
			"---------- WebSocket 启动程序成功 ----------\n",
		"---------- 服务器IP："+serverIp+" ----------\n",
		"---------- 服务器启动端口："+serverPort+" ----------\n",
		"---------- WebSocket端口："+webSocketPort+" ----------\n",
		"-------------------------------------------\n")

	_ = http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{
		// 允许跨域请求
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println("升级协议（允许跨域请求）", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
			return true
		},
	}).Upgrade(w, req, nil)

	if err != nil {
		http.NotFound(w, req)
		return
	}

	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	/*
		启用两个协程来处理
		客户端请求数据 client.read()
		和
		向客户端发送数据 client.write()
	*/
	go client.read()

	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
