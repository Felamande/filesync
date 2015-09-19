package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	svc "github.com/kardianos/service"

	"github.com/Felamande/filesync/syncer"
)

var run *bool = flag.Bool("run", false, "Run in the shell. -svcctl will be disabled.")
var controls *string = flag.String("svcctl", "install,start", "value:[start,stop,restart,install,uninstall], can be multiple values separated by commas")
var help *bool = flag.Bool("help", false, "Get help")

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if flag.Lookup("run") == nil && flag.Lookup("svcctl") == nil && flag.Lookup("help") == nil {
		flag.Usage()
		return
	}

	if *help {
		flag.Usage()
		return
	}

	p := &Program{}
	s, err := svc.New(p, &svc.Config{
		Name:        "Filesync",
		DisplayName: "FileSync Service",
		Description: "Filesync is a simple tool to sync files between multiple directory pairs.",
		Arguments:   []string{"-run"},
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if *run {
		err := s.Run()
		fmt.Println("-run with error: ", err.Error())
		return
	}

	Actions := strings.Split(*controls, ",")
	for _, action := range Actions {
		err := svc.Control(s, action)
		fmt.Println(err)
	}

}

func HelloNewPair() string {
	return "<!DOCTYPE html><head><script type='text/javascript' src='http://libs.baidu.com/jquery/2.0.3/jquery.min.js'></script></head><body>Hello</body>"
}

func NewPair(s *syncer.Syncer, r *http.Request) string {
	err := r.ParseForm()
	if err != nil {
		return err.Error()
	}

	lName := r.FormValue("left")
	rName := r.FormValue("right")
	err = s.NewPair(syncer.SyncConfig{true, true, true}, lName, rName)
	if err != nil {
		return err.Error()
	}
	return lName + " ==> " + rName

}
