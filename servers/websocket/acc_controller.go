/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 13:12
 */

package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gowebsocket/common"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"time"
)

// ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	fmt.Println("------------------------ acc_controller.go WebSocket Ping ------------------------")
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)
	data = "pong"
	return
}

// 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	// 解析发送过来的Json数据为Login格式
	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("------------------------ acc_controller.go WebSocket Login-ParameterIllegal ------------------------")
		fmt.Println("用户登录 解析数据失败", seq, err)
		return
	}
	// webSocket_request 用户登录 1577097036373-519199 ServiceToken
	fmt.Println("------------------------ acc_controller.go WebSocket Login ------------------------")
	fmt.Println("seq: ", seq, "ServiceToken: ", request.ServiceToken)

	if request.UserId == "" || len(request.UserId) > 32 {
		code = common.UnauthorizedUserId
		fmt.Println("------------------------ acc_controller.go WebSocket Login-UnauthorizedUserId ------------------------")
		fmt.Println("用户登录 非法的用户", seq, request.UserId)
		return
	}

	if !InAppIds(request.AppId) {
		code = common.Unauthorized
		fmt.Println("------------------------ acc_controller.go WebSocket Login-Unauthorized ------------------------")
		fmt.Println("用户登录 不支持的平台", seq, request.AppId)
		return
	}

	if client.IsLogin() {
		fmt.Println("------------------------ acc_controller.go WebSocket Login-HadLogin ------------------------")
		fmt.Println("用户登录 用户已经登录", client.AppId, client.UserId, seq)
		code = common.OperationFailure

		return
	}

	client.Login(request.AppId, request.UserId, currentTime)

	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("------------------------ acc_controller.go WebSocket Login-ServerError ------------------------")
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)
		return
	}

	// 用户登录
	login := &login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Client: client,
	}
	clientManager.Login <- login
	fmt.Println("------------------------ acc_controller.go WebSocket Login-Success ------------------------")
	fmt.Println("seq: ", seq, "addr: ", client.Addr, "userId: ",request.UserId)

	return
}

// 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat-Fail ------------------------")
		fmt.Println("心跳接口 解析数据失败", seq, err)
		return
	}
	fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat ------------------------")
	fmt.Println("appId: ", client.AppId, "userId: ", client.UserId)

	if !client.IsLogin() {
		fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat-UnLogin ------------------------")
		fmt.Println("心跳接口 用户未登录", client.AppId, client.UserId, seq)
		code = common.NotLoggedIn
		return
	}

	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat-unLogin ------------------------")
			fmt.Println("心跳接口 用户未登录", seq, client.AppId, client.UserId)
			return
		} else {
			code = common.ServerError
			fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat-GetUserOnlineInfo ------------------------")
			fmt.Println("seq: ", seq, "appId: ",client.AppId,"userId: ", client.UserId, "err: ",err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("------------------------ acc_controller.go WebSocket HeartBeat-SetUserOnlineInfo ------------------------")
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)
		return
	}

	return
}
