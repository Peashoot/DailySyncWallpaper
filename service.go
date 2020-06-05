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

	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows/registry"
)

func init() {
	// 获取配置
	config.getConf()
	// 防止重复启动（将pid保存到临时文件目录）
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

func main() {
	SetAutoRun()
	download()
	timer := cron.New()
	// 每0或30分时执行一次
	spec := config.CronRule
	timer.AddFunc(spec, func() {
		Println("cron running...")
		download()
	})
	timer.Start()
	select {}
}

// ProcExsit 判断进程是否启动(查询之前的pid文件是否存在，如果存在，查询pid程序是否在运行，运行程序是否和当前程序名称一致等)
func ProcExsit(tmpDir string) (err error) {
	iManPidFile, err := os.Open(tmpDir + "\\wallpaper.pid")
	defer iManPidFile.Close()

	if err == nil {
		filePid, err := ioutil.ReadAll(iManPidFile)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			proc, err := process.NewProcess(int32(pid))
			procCur, errCur := process.NewProcess(int32(os.Getpid()))
			if err == nil && errCur == nil {
				procName, err := proc.Name()
				procCurName, errCur := procCur.Name()
				if err == nil && errCur == nil && procName == procCurName {
					Println("DailySyncWallpaper is aleady launched.")
					return errors.New("[ERROR] DailySyncWallpaper is aleady launched")
				}
			}
		}
	}

	return nil
}

// SetAutoRun 设置开机自启
func SetAutoRun() {
	execpath := "\"" + os.Args[0] + "\""
	Println("Set executable path into auto run registry...")
	key, exists, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		Panicln(err)
	}
	defer key.Close()
	var oldvalue string
	if exists {
		Println("DailySyncWallpaper self-starting registry already exists.")
		oldvalue, _, _ = key.GetStringValue("DailySyncWallpaper")
	} else {
		Println("Create new registry of DailySyncWallpaper")
	}
	if oldvalue != execpath {
		Println("DailySyncWallpaper update path into self-starting registry.")
		key.SetStringValue("DailySyncWallpaper", execpath)
	}
	Println("Finish setting auto run registry...")
}
