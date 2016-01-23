package main

import (
    
	"flag"
	"fmt"
	"os"
	"strings"
    
    "github.com/Felamande/filesync/settings"
	"github.com/Felamande/filesync/log"
    // "github.com/Felamande/filesync/modules/utils"

	svc "github.com/kardianos/service"
)

var run = flag.Bool("run", false, "Run in the shell. -svcctl will be disabled.")
var controls = flag.String("svcctl", "install,start", "value:[start,stop,restart,install,uninstall], can be multiple values separated by commas")
var help = flag.Bool("help", false, "Get help")
var console = flag.Bool("console", false, "Print logs to the console instead of the log files.")

func main() {
    settings.Init()
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if flag.Lookup("run") == nil && flag.Lookup("svcctl") == nil && flag.Lookup("help") == nil && flag.Lookup("console") == nil {
		flag.Usage()
		return
	}

	if *help {
		flag.Usage()
		return
	}
    
	LogFile, err := os.OpenFile(settings.Log.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := &Program{
		Logger: log.New(LogFile, "[filesync]", log.Ldefault|log.Lmicroseconds),
	}

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
		if *console {
			p.Logger = log.New(os.Stdout, "[filesync]", log.Ldefault|log.Lmicroseconds)
			p.run()
			return
		}
		err := s.Run()
		fmt.Println("run with error: ", err)
		return
	}

	Actions := strings.Split(*controls, ",")
	for _, action := range Actions {
		err := svc.Control(s, action)
		fmt.Println(err)
	}

}