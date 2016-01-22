package main

import (
	// "bytes"
	"errors"
	"fmt"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
    "github.com/lunny/tango"
	"github.com/Felamande/filesync/log"
	"github.com/Felamande/filesync/syncer"
    "github.com/Felamande/filesync/server"
	svc "github.com/kardianos/service"
)

type Program struct {
	Syncer *syncer.Syncer
	Logger *log.Logger
	Config *syncer.SavedConfig
    Server *tango.Tango
	Folder string
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
    p.Server = server.Init(p.Syncer)
	go p.Syncer.Run(*p.Config)
    p.Server.Run(":9000")
}

//Stop Stop the program.
func (p *Program) Stop(s svc.Service) error {

	// c := syncer.SavedConfig{
	// 	Pairs: []syncer.SyncPairConfig{},
	// }
	// c.LogPath = p.Config.LogPath

	// for _, pair := range p.Syncer.SyncPairs {

	// 	c.Pairs = append(c.Pairs, syncer.SyncPairConfig{
	// 		Left:   pair.Left.Uri(),
	// 		Right:  pair.Right.Uri(),
	// 		Config: pair.Config,
	// 	})
	// }
	// c.Port = p.Config.Port
	// b, err := yaml.Marshal(&c)
	// if err != nil {
	// 	p.Logger.Warn(err.Error())
	// 	return err
	// }
	// rb, err := ioutil.ReadFile(filepath.Join(p.Folder, "config.yaml"))
	// if err != nil {
	// 	p.Logger.Warn(err.Error())
	// 	return err
	// }
	// if bytes.Equal(b, rb) {
	// 	p.Logger.Info("service stopped, config not been changed.")
	// 	return nil
	// }
	// err = ioutil.WriteFile(filepath.Join(p.Folder, "config.yaml"), b, 0777)
	p.Logger.Info("service stopped, config saved.")
	return nil
}

//ReadConfig config.yaml
func ReadConfig(ConfigFile string) (*syncer.SavedConfig, error) {

	config := &syncer.SavedConfig{
		Pairs: []syncer.SyncPairConfig{},
		Port:  20000,
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

	fmt.Println("Listen on port: ", config.Port)

	return config, nil

}
