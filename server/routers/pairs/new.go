package pairs

import (
	// "github.com/lunny/tango"
	// "errors"
    "github.com/Felamande/filesync/settings"
	"github.com/Felamande/filesync/server/models"
	"github.com/Felamande/filesync/server/routers/base"
	"github.com/Felamande/filesync/syncer"
	"github.com/tango-contrib/binding"
)

type NewPairRouter struct {
	base.BaseJSONRouter
	binding.Binder

	config settings.SyncConfig
    Form  models.NewPairForm
    Errors binding.Errors
}

func (r *NewPairRouter) Before() {
	r.BaseJSONRouter.Before()
	r.Form = models.NewPairForm{}
	if e := r.Bind(&r.Form); len(e) != 0 {
        r.Errors = e
	}
    
	if Formatter, ok := interface{}(r.Form).(models.Formatter); ok {
		Formatter.FormatTo(&r.config)
	}
}

func (r *NewPairRouter) Post() interface{} {
    if r.Errors.Len() !=0{
        r.JSON["err"] = r.Errors.ErrorMap()
        r.JSON["sucess"] =false
        return r.JSON
    }
    
    r.Logger.Info("New Pair:",r.config)
	err := syncer.Default().NewPair(r.config,r.Form.Left,r.Form.Right,r.Form.IgnoreRules)
    
    if err !=nil{
        r.JSON["sucess"] = false
        r.JSON["err"] = err.Error()   
    }else{
        r.JSON["sucess"] = true
        r.JSON["err"] = nil
    }
    
    
    
    return r.JSON
}
