package main

import (
	// "bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"path/filepath"
	"strconv"

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
	Logger *log.Logger
	Config *syncer.SavedConfig
	Folder string
}

//Start program
func (p *Program) Start(s svc.Service) error {
	go p.run()
	return nil
}

func (p *Program) run() {
	p.Logger.Info("start service.")
	p.Server.Map(p.Syncer)
	p.Server.Map(p.Logger)
	p.Server.Post("/new", NewPair)
	p.Server.Get("/new", HelloNewPair)
	syncer.SetLogger(p.Logger)
	go p.Syncer.Run(*p.Config)

	//	go func() {
	//		fmt.Println("watching config change")
	//		err := p.Syncer.WatchConfigChange(filepath.Join(p.Folder, "config.yaml"))
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	http.ListenAndServe(":"+strconv.Itoa(p.Config.Port), p.Server)
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
