/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-30
* Time: 12:27
 */

package websocket

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/grpcclient"
	"time"
)

// 查询所有用户
func UserList() (userList []string) {

	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)

		return
	}

	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList()
		} else {
			list, _ = grpcclient.GetUserList(server)
		}
		userList = append(userList, list...)
	}

	return
}

// 查询用户是否在线
func CheckUserOnline(appId uint32, userId string) (online bool) {
	// 全平台查询
	if appId == 0 {
		for _, appId := range GetAppIds() {
			online, _ = checkUserOnline(appId, userId)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appId, userId)
	}

	return
}

// 查询用户 是否在线
func checkUserOnline(appId uint32, userId string) (online bool, err error) {
	key := GetUserKey(appId, userId)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", appId, userId, err)

			return false, nil
		}

		fmt.Println("GetUserOnlineInfo", appId, userId, err)

		return
	}

	online = userOnline.IsOnline()

	return
}

// 给具体用户发送消息
func SendUserMessage(appId uint32, fromId string, toId string, msgId, message string) (sendResults bool, err error) {
	// GetTextMsgData
	data := models.GetTextMsgData(fromId, toId, msgId, message)
	// TODO::需要判断不在本机的情况
	sendResults, err = SendUserMessageLocal(appId, toId, data)
	if err != nil {
		fmt.Println("给用户发送消息", appId, fromId, toId, err)
	}
	return
}

// 给本机用户发送消息
func SendUserMessageLocal(appId uint32, toId string, data string) (sendResults bool, err error) {
	client := GetToUserClient(appId,toId)
	if client == nil {
		err = errors.New("用户不在线")
		return
	}
	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true
	return
}

// 给全体用户发消息
func SendUserMessageAll(appId uint32, fromId string, msgId, cmd, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		if IsLocal(server) {
			// GetTextMsgData
			data := models.GetMsgData(fromId, msgId, cmd, message)
			AllSendMessages(appId, fromId, data)
		} else {
			grpcclient.SendMsgAll(server, msgId, appId, fromId, cmd, message)
		}
	}
	return
}
