package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Felamande/filesync/log"
	"github.com/go-martini/martini"

	"os/signal"

	"github.com/Felamande/filesync/syncer"
)

func main() {
	s := syncer.New()
	m := martini.Classic()
	l := log.NewFileLogger("./.log/http.log")

	SigChan := make(chan os.Signal, 2)
	signal.Notify(SigChan, os.Interrupt, os.Kill)
	go handleSignal(SigChan, s)

	m.Map(s)
	m.Map(l)
	m.Post("/new", NewPair)
	m.Get("/new", HelloNewPair)
	go s.Run()
	http.ListenAndServe(":20000", m)

}

func NewPair(logger *log.FileLogger, s *syncer.Syncer, w http.ResponseWriter, r *http.Request) string {
	r.ParseForm()
	lName := r.FormValue("left")
	rName := r.FormValue("right")

	err := s.NewPair(syncer.SyncConfig{false, false, false}, lName, rName)

	if err != nil {
		return err.Error()
	}

	return "Success"

}

func HelloNewPair() string {
	return "<!DOCTYPE html><head><script type='text/javascript' src='http://libs.baidu.com/jquery/2.0.3/jquery.min.js'></script></head><body>Hello</body>"
}

func handleSignal(c chan os.Signal, s *syncer.Syncer) {
	for {
		select {
		case sig := <-c:
			fmt.Println(sig.String())
			c := syncer.SavedConfig{}
			for _, pair := range s.SyncPairs {
				fmt.Println(pair.Left.Uri())
				c.Pairs = append(c.Pairs, syncer.SyncPairConfig{
					Left:   pair.Left.Uri(),
					Right:  pair.Right.Uri(),
					Config: pair.Config,
				})
				b, _ := json.Marshal(c.Pairs)
				ioutil.WriteFile("config.json", b, 0777)
				os.Exit(0)
			}
		}
	}
}
