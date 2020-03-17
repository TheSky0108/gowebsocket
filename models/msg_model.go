/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-01
* Time: 10:40
 */

package models

import "gowebsocket/common"

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeVoice = "voice"
	MessageCmdMsg    = "msg"
	MessageCmdEnter  = "enter"
	MessageCmdExit   = "exit"

	SystemNoticeType  = "00"
	ChatNoticeType    = "01"
	PraiseNoticeType  = "02"
	DissNoticeType    = "02"
	CommentNoticeType = "03"
)

// 消息的定义
type Message struct {
	Target   string `json:"target"`   // 目标
	Type     string `json:"type"`     // 消息类型 text/image/
	Msg      string `json:"msg"`      // 消息内容
	From     string `json:"from"`     // 发送者
	To       string `json:"to"`       // 接收者
	SendTime uint64 `json:"sendTime"` // 发送时间（UNIX() 时间）
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

func NewTestMsgWithSendTime(from string, to string, Msg string, sendTime uint64) (message *Message) {

	message = &Message{
		Type:     MessageTypeText,
		From:     from,
		To:       to,
		Msg:      Msg,
		SendTime: sendTime,
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
		"sendTime":1582804913
		}
	}
}
*/
func getTextMsgData(cmd, fromId, toId, msgId, message string) string {
	textMsg := NewTestMsg(fromId, toId, message)
	head := NewResponseHead(msgId, cmd, common.OK, "Ok", textMsg)
	return head.String()
}

func getTextMsgDataWithSendTime(cmd, fromId, toId, msgId, message string, sendTime uint64) string {
	textMsg := NewTestMsgWithSendTime(fromId, toId, message, sendTime)
	head := NewResponseHead(msgId, cmd, common.OK, "Ok", textMsg)
	return head.String()
}

// 封装 Data
// 给 All 用户发送消息
func GetMsgData(userId, msgId, cmd, message string, sendTime uint64) string {
	return getTextMsgDataWithSendTime(cmd, userId, "All", msgId, message, sendTime)
}

// 封装 Data
// 给指定用户发送消息
func GetTextMsgData(fromId, toId, msgId, message string, sendTime uint64) string {
	return getTextMsgDataWithSendTime("msg", fromId, toId, msgId, message, sendTime)
}

// 用户进入消息
func GetTextMsgDataEnter(fromId, toId, msgId, message string) string {
	return getTextMsgData("enter", fromId, toId, msgId, message)
}

// 用户退出消息
func GetTextMsgDataExit(fromId, toId, msgId, message string) string {
	return getTextMsgData("exit", fromId, toId, msgId, message)
}
