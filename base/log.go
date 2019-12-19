package base

import (
	"fmt"
	"github.com/kakaisaname/props/kvs"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//Logrus格式化程序主要基于原始的Logrus.TextFormatter，
// 但有稍微修改的彩色输出和对日志条目前缀的支持，例如消息源后跟冒号。此外，还支持自定义颜色主题

var formatter *prefixed.TextFormatter
var lfh *utils.LineNumLogrusHook

//定义日志格式
//日志级别
//控制台高亮显示
//日志文件和滚动配置
func init() {
	//定义日志格式   文本格式																			***
	//prefixed 第三方包
	formatter = &prefixed.TextFormatter{}
	//设置高亮显示的色彩样式					强制高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false

	//强制格式化
	formatter.ForceFormatting = true
	//自定义高亮显示的颜色
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	//开启完整时间戳输出和时间戳格式 																	***
	formatter.FullTimestamp = true
	//设置时间格式
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	//设置日志formatter																				***
	log.SetFormatter(formatter)

	log.SetOutput(colorable.NewColorableStdout())
	//日志级别，通过环境变量来设置
	// 后期可以变更到配置中来设置

	//这个需要自己去配置  日志级别。通过系统环境变量来设置 											   ***
	if os.Getenv("log.debug") == "true" {
		log.SetLevel(log.DebugLevel)
	}
	//开启调用函数、文件、代码行信息的输出
	log.SetReportCaller(true)
	//设置函数、文件、代码行信息的输出的hook
	SetLineNumLogrusHook()
}

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	log.AddHook(lfh)
}

//将滚动日志writer共享给iris glog output
var log_writer io.Writer

//初始化log配置，配置logrus日志文件滚动生成和 														***
func InitLog(conf kvs.ConfigSource) {
	//设置日志输出级别
	level, err := log.ParseLevel(conf.GetDefault("log.level", "info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	if conf.GetBoolDefault("log.enableLineLog", true) {
		lfh.EnableFileNameLog = true
		lfh.EnableFuncNameLog = true
	} else {
		lfh.EnableFileNameLog = false
		lfh.EnableFuncNameLog = false
	}

	//配置日志输出目录
	logDir := conf.GetDefault("log.dir", "./logs")
	logTestDir, err := conf.Get("log.test.dir")
	if err == nil {
		logDir = logTestDir
	}
	logPath := logDir //+ "/logs"
	logFilePath, _ := filepath.Abs(logPath)
	log.Infof("log dir: %s", logFilePath)
	logFileName := conf.GetDefault("log.file.name", "red-envelop")
	maxAge := conf.GetDurationDefault("log.max.age", time.Hour*24)
	rotationTime := conf.GetDurationDefault("log.rotation.time", time.Hour*1)
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	//设置滚动日志输出writer
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	//设置日志文件输出的日志格式
	formatter := &log.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, formatter)

	log.AddHook(lfHook)
	//
	log_writer = writer

}
