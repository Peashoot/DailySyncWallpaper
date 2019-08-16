// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/cron"
	"github.com/kardianos/service"
	"golang.org/x/sys/windows/registry"
)

func init() {
	config.getConf()
	iManPid := fmt.Sprint(os.Getpid())
	tmpDir := os.TempDir()
	if err := ProcExsit(tmpDir); err == nil {
		pidFile, _ := os.Create(tmpDir + "\\wallpaper.pid")
		defer pidFile.Close()
		pidFile.WriteString(iManPid)
	} else {
		os.Exit(1)
	}
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	SetAutoRun()
	Println(fmt.Sprintf("I'm running %v.", service.Platform()))
	download()
	timer := cron.New()
	spec := "0 */30 * * * ?"
	timer.AddFunc(spec, func() {
		Println("cron running...")
		download()
	})
	timer.Start()
	select {}
}

// 判断进程是否启动
func ProcExsit(tmpDir string) (err error) {
	iManPidFile, err := os.Open(tmpDir + "\\wallpaper.pid")
	defer iManPidFile.Close()

	if err == nil {
		filePid, err := ioutil.ReadAll(iManPidFile)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				return errors.New("[ERROR] DailySyncWallpaper已启动.")
			}
		}
	}

	return nil
}

func SetAutoRun() {
	execpath := "\"" + os.Args[0] + "\""
	Println("Set executable path into auto run registry...")
	key, exists, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		Fatalln(err)
	}
	defer key.Close()
	var oldvalue string
	if exists {
		Println("DailySyncWallpaper键已存在")
		oldvalue, _, _ = key.GetStringValue("DailySyncWallpaper")
	} else {
		Println("新建注册表键DailySyncWallpaper")
	}
	if oldvalue != execpath {
		key.SetStringValue("DailySyncWallpaper", execpath)
	}
	Println("Success set auto run registry...")
}
