package glog

import (
	"log"
	"os"
	"fmt"
	runtimeDebug "runtime/debug"
)



//定义3个log变量,用来输出不同的内容
var info log.Logger = *log.New(nil, "[Glog-info]", log.LstdFlags)
var error log.Logger = *log.New(nil, "[Glog-error]", log.LstdFlags)
var debug log.Logger = *log.New(nil, "[Glog-debug]", log.LstdFlags)

var iFile = new(os.File)
var dFile = new(os.File)
var eFile = new(os.File)


//三个级别的日志
//使用包名调用,简单粗暴
func Info(str string) {
	fmt.Println(str)
	info.Println(str)
}
func Debug(str string) {
	fmt.Println(str)
	debug.Println(str)
}
func Error(str string) {

	//打印堆栈
	stack := runtimeDebug.Stack()
	info := str + "\n" + string(stack)
	fmt.Println(info)
	error.Println(info)

}




//初始化info日志的writer
func initInfo() *os.File {
	file, err := os.OpenFile("/Users/yehao/my/info.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		panic("info日志初始化错误")
	}
	iFile = file
	return file

}


//初始化debug日志的writer
func initDebug() *os.File {
	file, err := os.OpenFile("/Users/yehao/my/debug.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		panic("debug日志初始化错误")
	}
	dFile = file
	return file

}



//初始化error日志的writer
func initError() *os.File {
	file, err := os.OpenFile("/Users/yehao/my/error.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		panic("error日志初始化错误")
	}
	eFile = file
	return file

}

//初始化日志目录
func init() {

	fmt.Println("初始化glog...")

	info.SetOutput(initInfo())
	error.SetOutput(initError())
	debug.SetOutput(initDebug())
}

//依次关闭文件
//用deffer防止 异常导致其他的file没有关闭
func Close() {
	defer func() {
		iFile.Close()
		fmt.Println("info日志文件关闭成功!")
		fmt.Println("所有日志文件关闭成功")
	}()
	defer func() {
		dFile.Close()
		fmt.Println("debug日志文件关闭成功!")
	}()
	defer func() {
		eFile.Close()
		fmt.Println("error日志文件关闭成功!")
	}()
}