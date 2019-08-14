// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"flag"
	"fmt"

	"github.com/cron"

	"github.com/kardianos/service"
)

func init() {
	config.getConf()
	LoadLogger()
}

// Program structures.
// Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		Println("Running in terminal.")
	} else {
		Println("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	Println(fmt.Sprintf("I'm running %v.", service.Platform()))
	download()
	timer := cron.New()
	spec := "0 */30 * * * ?"
	timer.AddFunc(spec, func() {
		Println("cron running...")
		download()
	})
	timer.Start()
	for {
		select {
		case <-p.exit:
			timer.Stop()
			return nil
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	Println("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "DailySyncWallpaper",
		DisplayName: "必应壁纸同步服务",
		Description: "定时同步bing壁纸",
		Option:      options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		Fatalln(err)
	}
	errs := make(chan error, 5)
	_, err = s.Logger(errs)
	if err != nil {
		Fatalln(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				Fatalln(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			Println(fmt.Sprintf("Valid actions: %q", service.ControlAction))
			Fatalln(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		Fatalln(err)
	}
}
