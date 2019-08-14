package main

import (
	"log"

	"github.com/kardianos/service"
)

var logger service.Logger

//写日志
func Println(msg string) {
	logger.Info(msg)
	log.Println(msg)
}

//写错误日志
func Fatalln(err error) {
	logger.Error(err)
	log.Fatalln(err)
}
