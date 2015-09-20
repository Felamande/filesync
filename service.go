package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/Felamande/filesync/log"
)
import (
	"github.com/Felamande/filesync/syncer"
	"github.com/go-martini/martini"
	svc "github.com/kardianos/service"
)

type Program struct {
	Syncer *syncer.Syncer
	Server *martini.ClassicMartini
	Logger *log.FileLogger
	Config *syncer.SavedConfig
	Folder string
}

func (p *Program) Start(s svc.Service) error {
	go p.run()
	return nil
}

func (p *Program) run() {

	p.Server.Map(p.Syncer)
	p.Server.Map(p.Logger)
	p.Server.Post("/new", NewPair)
	p.Server.Get("/new", HelloNewPair)
	go p.Syncer.Run(*p.Config)
	//	go func() {
	//		fmt.Println("watching config change")
	//		err := p.Syncer.WatchConfigChange(filepath.Join(p.Folder, "config.yaml"))
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	http.ListenAndServe(":"+p.Config.Port, p.Server)
}

func (p *Program) Stop(s svc.Service) error {
	c := syncer.SavedConfig{
		Pairs: []syncer.SyncPairConfig{},
	}

	for _, pair := range p.Syncer.SyncPairs {

		c.Pairs = append(c.Pairs, syncer.SyncPairConfig{
			Left:   pair.Left.Uri(),
			Right:  pair.Right.Uri(),
			Config: pair.Config,
		})
	}
	c.Port = p.Config.Port
	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(p.Folder, "config.yaml"), b, 0777)
	return err
}

func ReadConfig(ConfigFile string) (*syncer.SavedConfig, error) {

	var config *syncer.SavedConfig = &syncer.SavedConfig{
		Pairs: []syncer.SyncPairConfig{},
		Port:  "20000",
	}

	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if len(config.Port) == 0 {
		config.Port = "20000"
		fmt.Println("Listen on the default port: 2000")
	}
	fmt.Println("Listen on port: " + config.Port)

	return config, nil

}
