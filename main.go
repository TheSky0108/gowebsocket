/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 09:59
 */

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gowebsocket/lib/redislib"
	"gowebsocket/middleware"
	"gowebsocket/routers"
	"gowebsocket/servers/grpcserver"
	"gowebsocket/servers/task"
	"gowebsocket/servers/websocket"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	// 初始化配置文件
	initConfig()

	// 初始化路由(包含日志)
	router := initRouter()

	// 初始化定时任务
	initTimerTask()

	// WS服务启动
	go websocket.StartWebSocket()

	// grpc
	go grpcserver.Init()

	// 体验地址
	go openDemoUrl()

	// Http服务启动
	httpPort := viper.GetString("app.httpPort")
	_ = http.ListenAndServe(":"+httpPort, router)

}

func initConfig() {
	// https://github.com/spf13/viper
	viper.SetConfigName("./config/app") // name of config file (without extension)
	viper.AddConfigPath(".")            //
	// viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	// viper.AddConfigPath(".")               // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 初始化Redis相关配置
	initRedis()

	fmt.Println("---------------- main.go 配置文件 ----------------")
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))
	fmt.Println("-------------------------------------------------")
}

func initRedis() {
	redislib.ExampleNewClient()
}

func initRouter() *gin.Engine {

	// 初始化gin启动模式，默认DebugMode
	//initGinMode(gin.ReleaseMode)

	fmt.Println("---------------- main.go 初始化路由 --------------")
	// 如果想完全使用自定义的Logger() 则需要使用gin.New()来生成对象
	// 然后再手动指定Logger(),Recovery()
	// 否则gin会使用两个日志框架，一个默认的，一个自定义的，影响性能
	var router *gin.Engine

	//使用自定义日志框架
	myLogger := middleware.LoggerToFile()
	if myLogger != nil {
		router = gin.New()
		router.Use(myLogger, gin.Recovery())
	} else {
		router = gin.Default()
	}

	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()
	fmt.Println("路由初始化完成")
	fmt.Println("-------------------------------------------------")
	return router
}

func initGinMode(ginMode string) {
	gin.SetMode(ginMode)

	// gin.SetMode() 要在 gin.Default() 之前调用
	//gin.SetMode(gin.DebugMode) 	// 默认启动模式
	//gin.SetMode(gin.ReleaseMode)  // 发布模式
	//gin.SetMode(gin.TestMode)     // 测试模式

	/* 如果使用 gin.DebugMode 则会在控制台打印如下内容
	[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

	[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	 - using env:	export GIN_MODE=release
	 - using code:	gin.SetMode(gin.ReleaseMode)

	[GIN-debug] Loaded HTML Templates (3):
		-
		- index.html

	[GIN-debug] GET    /user/list                --> gowebsocket/controllers/user.List (3 handlers)
	[GIN-debug] GET    /user/online              --> gowebsocket/controllers/user.Online (3 handlers)
	[GIN-debug] POST   /user/sendMessage         --> gowebsocket/controllers/user.SendMessage (3 handlers)
	[GIN-debug] POST   /user/sendMessageAll      --> gowebsocket/controllers/user.SendMessageAll (3 handlers)
	[GIN-debug] GET    /system/state             --> gowebsocket/controllers/systems.Status (3 handlers)
	[GIN-debug] GET    /home/index               --> gowebsocket/controllers/home.Index (3 handlers)
	*/
}

func initTimerTask() {
	// 定时任务
	task.Init()
	// 服务注册
	task.ServerInit()
}

func openDemoUrl() {
	time.Sleep(1000 * time.Millisecond)
	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"
	fmt.Println("---------------- main.go ----------------")
	fmt.Println("访问页面体验:", httpUrl)
	fmt.Println("-----------------------------------------")
	cmd := exec.Command("openDemoUrl", httpUrl)
	_, _ = cmd.Output()
}
