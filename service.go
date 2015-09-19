package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Felamande/filesync/log"
	"github.com/kardianos/osext"
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
	folder, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := ReadConfig(filepath.Join(folder, "config.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)

	p = &Program{
		Config: config,
		Syncer: syncer.New(),
		Server: martini.Classic(),
		Logger: log.NewFileLogger(filepath.Join(folder, "./.log/service.log")),
		Folder: folder,
	}
	p.Server.Map(p.Syncer)
	p.Server.Map(p.Logger)
	p.Server.Post("/new", NewPair)
	p.Server.Get("/new", HelloNewPair)
	go p.Syncer.Run(*p.Config)
	http.ListenAndServe(p.Config.Port, p.Server)
}

func (p *Program) Stop(s svc.Service) error {
	c := syncer.SavedConfig{}
	for _, pair := range p.Syncer.SyncPairs {

		c.Pairs = append(c.Pairs, syncer.SyncPairConfig{
			Left:   pair.Left.Uri(),
			Right:  pair.Right.Uri(),
			Config: pair.Config,
		})
	}
	c.Port = p.Config.Port
	b, err := json.Marshal(&c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(p.Folder, "config.json"), b, 0777)
	return err
}

func ReadConfig(ConfigFile string) (*syncer.SavedConfig, error) {

	var config *syncer.SavedConfig = &syncer.SavedConfig{
		Pairs: []syncer.SyncPairConfig{},
		Port:  ":20000",
	}

	configFile, err := os.Open(ConfigFile)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	d := json.NewDecoder(configFile)
	err = d.Decode(config)

	if err != nil {
		return nil, err
	}

	if len(config.Port) == 0 {
		config.Port = ":20000"
	}
	if config.Port[0] != ':' {
		config.Port = ":" + config.Port
	}

	return config, nil

}
