package server

import (
    "github.com/Felamande/filesync/syncer"
    "github.com/lunny/tango"
    
    //middlewares 
    "github.com/Felamande/filesync/server/middlewares/time"
    "github.com/tango-contrib/events"
    "github.com/tango-contrib/renders"
    
    //routers
    "github.com/Felamande/filesync/server/routers/pairs"
    "github.com/Felamande/filesync/server/routers/page"
    
    
    
    
)

var t *tango.Tango

func init(){
    t = tango.New()
}

func Init(s *syncer.Syncer) *tango.Tango{
    t.Use(new(time.TimeHandler))
    t.Use(events.Events())
    t.Use(tango.ClassicHandlers...)
    t.Use(renders.New(renders.Options{
        Reload:true,
        Directory:"./server/templates",
        Charset:"UTF-8",
        DelimsLeft:"{%",
        DelimsRight:"%}",
    }))
    
    t.Get("/pair/all",new(pairs.NewPairRouter))
    t.Get("/",new(page.HomeRouter))
    return t
}
