package logs

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var(
	MyLogger *myLogger
)

type myLogger struct {
	*logrus.Logger
	File *os.File
}

type myHook struct {
	FileName string
	Line string
	Skip int
	levels []logrus.Level
}

//实现 logrus.Hook 接口
func (hook *myHook)Fire(entry *logrus.Entry) error  {
	fileName, line := findCaller(hook.Skip)
	entry.Data[hook.FileName] = fileName
	entry.Data[hook.Line] = line
	return nil
}

//实现 logrus.Hook 接口
func (hook *myHook)Levels() []logrus.Level  {
	return logrus.AllLevels
}

//自定义hook
func NewMyHook(levels ...logrus.Level) logrus.Hook  {
	hook := myHook{
		FileName:"filePath",
		Line:"line",
		Skip:5,
		levels:levels,
	}
	if len(hook.levels) == 0{
		hook.levels = logrus.AllLevels
	}
	return &hook
}

func getCaller(skip int) (string, int)  {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println("getCaller", pc, file, line, ok)
	if !ok{
		return "", 0
	}
	n := 0
	//获取执行代码的文件名
	for i:=len(file) -1; i >0; i--{
		if string(file[i]) == "/" {
			n++
			if n >= 2 {
				//fmt.Println(n >= 2, file)
				file = file[i+1:]
				break
			}
		}
	}
	return file,line
}

func findCaller(skip int) (string,int)  {
	file := ""
	line := 0
	for i:=0; i<10; i++{
		file, line = getCaller(skip+i)
		//fmt.Println("findCaller", file, line)
		//文件名不能以logrus开头
		if !strings.HasPrefix(file, "logrus"){
			break
		}
	}
	return file, line
}

//自定义logger
func NewLogger(level logrus.Level, format logrus.Formatter, hook logrus.Hook) *logrus.Logger {
	log := logrus.New()
	log.Level = level
	log.SetFormatter(format)
	log.Hooks.Add(hook)
	return log
}

//初始化配置
func init() {
	path := "/var/log/kylin-ccm/ccm.log"

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	mw := io.MultiWriter(os.Stdout, file)
	if err == nil {
		logrus.SetOutput(mw)
	} else {
		logrus.Info("failed to writer log to file.")
	}

	MyLogger = &myLogger{
		File:file,
	}
	MyLogger.Logger = NewLogger(logrus.DebugLevel, &logrus.TextFormatter{FullTimestamp:true, TimestampFormat: "15:04:05 MST"}, NewMyHook())
	MyLogger.Logger.Out = MyLogger.File
}