package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	logFilePath := viper.GetString("app.logFilePath")
	logFileName := viper.GetString("app.logFileName")
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	src, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModeAppend)
	if err != nil {
		// 如果创建失败，尝试手动创建一下
		fmt.Println("**********日志文件创建失败**********")
		fmt.Println(src, err)
		return nil
	}
	fmt.Println("**********日志文件创建成功**********")
	// 实例化
	logger := logrus.New()

	// 设置输出
	logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	// 设置日志格式-文本格式
	// 单独格式化日期
	//logger.SetFormatter(&logrus.TextFormatter{
	//	TimestampFormat:"2006-01-02 15:04:05",
	//})

	// 设置日志格式-Json格式
	// 单独格式化日期
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		executionTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志输出格式
		// 输出到文件
		logger.WithFields(logrus.Fields{
			"req_method":     reqMethod,
			"req_uri":        reqUri,
			"status_code":    statusCode,
			"execution_time": executionTime,
			"client_ip":      clientIP,
		}).Info()

		// 输出到控制台
		fmt.Println("------------------------ Request ------------------------")
		fmt.Printf(
			"req_method: %s | req_uri: %s | status_code: %d | execution_time: %s | client_ip: %s\n",
			reqMethod, reqUri, statusCode, executionTime, clientIP)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
