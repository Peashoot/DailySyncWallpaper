package main

import (
	"log"
	"os"
)

// 日志操作类
var localLog *log.Logger

// 创建本地日志文件
func LoadLogger() {
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		log.Fatalln("fail to create service.log file!")
	}
	//defer file.Close()
	localLog = log.New(file, "", log.LstdFlags|log.Lshortfile) // 日志文件格式:log包含时间及文件行数
}

//写日志
func Println(msg string) {
	localLog.Println(msg)
	log.Println(msg)
}

//写错误日志
func Fatalln(err error) {
	localLog.Fatalln(err)
	log.Fatalln(err)
}
