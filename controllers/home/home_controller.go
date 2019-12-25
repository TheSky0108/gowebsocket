/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package home

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// 查看用户是否在线
func Index(c *gin.Context) {
	//给前端页面传递数据
	data := gin.H{
		"title": viper.GetString("html.title"),
		"contentTypeForm": viper.GetString("html.contentTypeForm"),
		"sendMessageUrl": viper.GetString("html.sendMessageUrl"),
		"sendMessageAllUrl": viper.GetString("html.sendMessageAllUrl"),
		"getUserListUrl": viper.GetString("html.getUserListUrl"),
		"webSocketUrl": viper.GetString("html.webSocketUrl"),
	}
	c.HTML(http.StatusOK, "index.html", data)
}
