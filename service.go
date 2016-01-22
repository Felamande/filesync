package main

import (
	// "bytes"

	"github.com/Felamande/filesync/log"
	"github.com/Felamande/filesync/syncer"
    "github.com/Felamande/filesync/server"
	svc "github.com/kardianos/service"
)

type Program struct {
	Logger *log.Logger
}

//Start program
func (p *Program) Start(s svc.Service) error {
	go p.run()
	return nil
}

func (p *Program) run() {
	p.Logger.Info("start service.")
	syncer.SetLogger(p.Logger)

	go syncer.Default().Run()
    server.Run()
}

//Stop Stop the program.
func (p *Program) Stop(s svc.Service) error {
	p.Logger.Info("service stopped, config saved.")
	return nil
}

//ReadConfig config.yaml