package main

import (
	// "bytes"
	"errors"
	"io/ioutil"
    "github.com/Felamande/filesync/settings"
    // "github.com/Felamande/filesync/server/modules/utils"
	yaml "gopkg.in/yaml.v2"
	"github.com/Felamande/filesync/log"
	"github.com/Felamande/filesync/syncer"
    "github.com/Felamande/filesync/server"
	svc "github.com/kardianos/service"
)

type Program struct {
	Syncer *syncer.Syncer
	Logger *log.Logger
	Config *syncer.SavedConfig
}

//Start program
func (p *Program) Start(s svc.Service) error {
	go p.run()
	return nil
}

func (p *Program) run() {
	p.Logger.Info("start service.")
	syncer.SetLogger(p.Logger)
    if p.Syncer == nil{
        panic("nil")
    }
	go p.Syncer.Run(*p.Config)
    server.Run(settings.Port)
}

//Stop Stop the program.
func (p *Program) Stop(s svc.Service) error {
	p.Logger.Info("service stopped, config saved.")
	return nil
}

//ReadConfig config.yaml
func ReadConfig(ConfigFile string) (*syncer.SavedConfig, error) {

	config := &syncer.SavedConfig{
		Pairs: []syncer.SyncPairConfig{},
	}
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if len(config.LogPath) == 0 {
		return nil, errors.New("Need a log path in config.yaml")
	}

	return config, nil

}
