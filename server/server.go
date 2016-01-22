package server

import (
	"github.com/Felamande/filesync/server/modules/utils"
	"github.com/Felamande/filesync/settings"
	"github.com/lunny/tango"
	//middlewares
	"github.com/Felamande/filesync/server/middlewares/time"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"

	//routers
	"github.com/Felamande/filesync/server/routers/page"
	"github.com/Felamande/filesync/server/routers/pairs"
)

var t *tango.Tango

func Run(addr interface{}) {

	t = tango.New()

	t.Use(tango.Static(tango.StaticOptions{
		RootPath: settings.LocalStatic,
	}))
	t.Use(new(time.TimeHandler))
	t.Use(events.Events())
	t.Use(tango.ClassicHandlers...)
	t.Use(renders.New(renders.Options{
		Reload:      true,
		Directory:   settings.TplHome,
		Charset:     settings.TplCharset,
		DelimsLeft:  settings.DelimesLeft,
		DelimsRight: settings.DelimesRight,
		Funcs:       utils.DefaultFuncs(),
	}))

	t.Get("/pair/all", new(pairs.NewPairRouter))
	t.Get("/", new(page.HomeRouter))

	t.Run(addr)
}
