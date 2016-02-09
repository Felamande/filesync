package page

import (
    // "net/http"
	"github.com/tango-contrib/renders"
	"github.com/Felamande/filesync/server/routers/base"
)

type HomeRouter struct {
	base.BaseTplRouter
}

func (r *HomeRouter) Get() {
    if r.Data == nil{
        r.Data = make(renders.T)
    }
	r.Data["title"] = "filesync dashboard "
    r.Tpl = "home.html"
    
    r.Render(r.Tpl,r.Data)
}
