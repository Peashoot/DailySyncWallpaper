package main

import (
	"log"
	"os"
)

var localLog *log.Logger

func LoadLogger() {
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		logger.Error(err)
	}
	localLog = log.New(file, "", log.LstdFlags)
}

//写日志
func Println(msg string) {
	localLog.Println(msg)
	log.Println(msg)
}

//写错误日志
func Fatalln(err error) {
	localLog.Println(err)
	log.Fatalln(err)
}
