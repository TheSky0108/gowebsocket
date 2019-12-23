/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-01
* Time: 10:40
 */

package models

import "gowebsocket/common"

const (
	MessageTypeText = "text"

	MessageCmdMsg   = "msg"
	MessageCmdEnter = "enter"
	MessageCmdExit  = "exit"
)

// 消息的定义
type Message struct {
	Target string `json:"target"` // 目标
	Type   string `json:"type"`   // 消息类型 text/img/
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
	To     string `json:"to"`     // 接收者
}

func NewTestMsg(from string, to string, Msg string) (message *Message) {

	message = &Message{
		Type: MessageTypeText,
		From: from,
		To:   to,
		Msg:  Msg,
	}

	return
}

/*
seq: 消息的唯一id
cmd: 消息的动作：enter(login)、msg、exit(logout)
response：返回前端的数据
from： 消息发送者id
to：   消息接收者id
{
"seq":"1577033806984-838426",
"cmd":"msg",
"response":{
	"code":200,
	"codeMsg":"Ok",
	"data":{
		"target":"",
		"type":"text",
		"msg":"213123213213",
		"from":"29868",
     	"to":"29868"
		}
	}
}
*/
func getTextMsgData(cmd, fromId, toId, msgId, message string) string {
	textMsg := NewTestMsg(fromId, toId, message)
	head := NewResponseHead(msgId, cmd, common.OK, "Ok", textMsg)
	return head.String()
}

// 给全体用户发消息，文本消息
// userId, msgId, cmd, message
func GetMsgData(userId, msgId, cmd, message string) string {
	return getTextMsgData(cmd, userId, "All", msgId, message)
}

// 给具体用户发消息，文本消息
// toId, msgId, message
func GetTextMsgData(fromId, toId, msgId, message string) string {
	return getTextMsgData("msg", fromId, toId, msgId, message)
}

// 用户进入消息
func GetTextMsgDataEnter(fromId, toId, msgId, message string) string {
	return getTextMsgData("enter", fromId, toId, msgId, message)
}

// 用户退出消息
func GetTextMsgDataExit(fromId, toId, msgId, message string) string {
	return getTextMsgData("exit", fromId, toId, msgId, message)
}
