/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 16:24
 */

package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"runtime/debug"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// 用户登录
type login struct {
	AppId  uint32
	UserId string
	Client *Client
}

// 读取客户端数据
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppId, l.UserId)

	return
}

// 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppId         uint32          // 登录的平台Id app/web/ios
	UserId        string          // 用户Id，用户登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
}

// 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// 读取客户端数据
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppId, c.UserId)
	return
}

// 读取客户端数据的异步处理程序
func (c *Client) read() {
	// 防止发生程序崩溃，所以需要捕获异常
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		fmt.Println("读取客户端数据 关闭send", c)
		// read() 读取失败时，close(c.Send) 关闭 write()
		close(c.Send)
	}()

	// 循环读取客户端发送的数据并处理
	// 如果读取数据失败了，return会触发defer关闭channel
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端数据 错误", c.Addr, err)
			return
		}

		// 处理程序
		fmt.Println("读取客户端数据 处理:", string(message))
		ProcessData(c, message)
	}
}

// 向客户端写数据的异步处理程序
func (c *Client) write() {
	// 防止发生程序崩溃，所以需要捕获异常
	defer func() {
		if r := recover(); r != nil {
			// 为了显示异常崩溃位置这里使用 string(debug.Stack()) 打印调用堆栈信息
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		clientManager.Unregister <- c
		// write()写入失败时，c.Socket.Close()关闭 read()
		_ = c.Socket.Close()
		fmt.Println("Client发送数据 defer", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 如果写入数据失败了，可能连接有问题，就return，触发defer去关闭连接
				fmt.Println("Client发送数据 关闭连接", c.Addr, "ok", ok)
				return
			}
			// 向客户端写数据
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// 读取客户端数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

// 读取客户端数据
func (c *Client) close() {
	close(c.Send)
}

// 用户登录
func (c *Client) Login(appId uint32, userId string, loginTime uint64) {
	c.AppId = appId
	c.UserId = userId
	c.LoginTime = loginTime
	// 登录成功=心跳一次
	c.Heartbeat(loginTime)
}

// 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	// heartbeatExpirationTime = 6 * 60 即 6 min
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// 是否登录了
func (c *Client) IsLogin() (isLogin bool) {

	// 用户登录了
	if c.UserId != "" {
		isLogin = true

		return
	}

	return
}
