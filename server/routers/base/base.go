package base

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type BaseJSONRouter struct {
	tango.Log
	tango.Json
    JSON  map[string]interface{}
}

func (r *BaseJSONRouter) Before() {
	r.JSON = make(map[string]interface{})
}

func (r *BaseJSONRouter) After() {

}

type BaseTplRouter struct {
	tango.Ctx
	tango.Log
	renders.Renderer
	Tpl  string
	Data renders.T
}

func (r *BaseTplRouter) Before() {
	r.Data = make(renders.T)
	r.Data["msg"] = "Hello, welcome to dashboard!"
}

func (r *BaseTplRouter) After() {
    
}
