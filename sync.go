package main

import (
	"net/http"

	"github.com/Felamande/filesync/log"
	"github.com/go-martini/martini"

	"github.com/Felamande/filesync/syncer"
)

func main() {
	s := syncer.New()
	m := martini.Classic()
	l := log.NewFileLogger("./.log/http.log")
	m.Map(s)
	m.Map(l)
	m.Post("/new", NewPair)
	m.Get("/new", HelloNewPair)
	go s.Run()
	http.ListenAndServe(":20000", m)

}

func NewPair(logger *log.FileLogger, s *syncer.Syncer, w http.ResponseWriter, r *http.Request) int {
	r.ParseForm()
	_, lExist := r.Form["left"]
	_, rExist := r.Form["right"]

	if !lExist || !rExist {
		return 400

	}
	return 200

}

func HelloNewPair() string {
	return "Hello"
}
