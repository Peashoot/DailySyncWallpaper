package main

import (
	"log"
	"os"
)

//写日志
func Println(msg string) {
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	localLog := log.New(file, "", log.LstdFlags)
	localLog.Println(msg)
	log.Println(msg)
}

//写错误日志
func Fatalln(err error) {
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE, 777)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	localLog := log.New(file, "", log.LstdFlags)
	localLog.Println(err)
	log.Fatalln(err)
}
