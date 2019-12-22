/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gowebsocket/common"
	"gowebsocket/controllers"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/websocket"
	"strconv"
)

// 查看全部在线用户
func List(c *gin.Context) {

	appIdStr := c.Query("appId")
	appId, _ := strconv.ParseInt(appIdStr, 10, 32)

	fmt.Println("http_request 查看全部在线用户", appId)

	data := make(map[string]interface{})

	userList := websocket.UserList()
	data["userList"] = userList

	controllers.Response(c, common.OK, "", data)
}

// 查看用户是否在线
func Online(c *gin.Context) {

	userId := c.Query("userId")
	appIdStr := c.Query("appId")

	fmt.Println("http_request 查看用户是否在线", userId, appIdStr)
	appId, _ := strconv.ParseInt(appIdStr, 10, 32)

	data := make(map[string]interface{})

	online := websocket.CheckUserOnline(uint32(appId), userId)
	data["userId"] = userId
	data["online"] = online

	controllers.Response(c, common.OK, "", data)
}

// 给用户发送消息
func SendMessage(c *gin.Context) {
	// 获取参数
	appIdStr := c.PostForm("appId") // appIds 一个用户在多个平台登录,比如：101-Web 102-iOS 103-Android
	fromId := c.PostForm("fromId")
	toId := c.PostForm("toId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")

	fmt.Println(
		"------------------------\n"+
			"http请求SendMessage\n",
		"appId: "+appIdStr+"\n",
		"发送者id: "+fromId+"\n",
		"接收者id: "+toId+"\n",
		"消息id: "+msgId+"\n",
		"消息内容: "+message+"\n",
		"------------------------\n")
	appId, _ := strconv.ParseInt(appIdStr, 10, 32)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgId) {
		fmt.Println("给具体用户发送消息，msgId重复提交:", msgId)
		controllers.Response(c, common.OK, "", data)
		return
	}
	//-------以上代码和 SendMessageAll() 完全一样
	sendResults, err := websocket.SendUserMessage(uint32(appId), fromId, toId, msgId, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)
}

// 给全员发送消息
func SendMessageAll(c *gin.Context) {
	// 获取参数
	appIdStr := c.PostForm("appId")
	fromId := c.PostForm("fromId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")

	fmt.Println(
		"------------------------\n"+
			"http请求SendMessageAll\n",
		"appId: "+appIdStr+"\n",
		"发送者id: "+fromId+"\n",
		"消息id: "+msgId+"\n",
		"消息内容: "+message+"\n",
		"------------------------\n")
	appId, _ := strconv.ParseInt(appIdStr, 10, 32)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgId) {
		fmt.Println("给全体用户发送消息，msgId重复提交:", msgId)
		controllers.Response(c, common.OK, "", data)
		return
	}
	//-------以上代码和 SendMessage() 完全一样
	sendResults, err := websocket.SendUserMessageAll(uint32(appId), fromId, msgId, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)
}
