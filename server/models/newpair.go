package models

import (
	"net/http"

	"github.com/Felamande/filesync/settings"
	"github.com/tango-contrib/binding"
)

type NewPairForm struct {
	Left          string   `form:"left" binding:"Required"`
	Right         string   `form:"right" binding:"Required"`
	SyncDel       string   `form:"syncdel"`
	SyncRename    string   `form:"syncrename"`
	CoverSameName string   `form:"cover"`
	IgnoreRules   []string `form:"ignore"`
}

func (*NewPairForm) Validate(r *http.Request, e binding.Errors) binding.Errors {

	var err error

	err = checkRadioParam(r, "syncdel")
	err = checkRadioParam(r, "syncrename")
	err = checkRadioParam(r, "coversamename")
	if err == nil {
		return e
	}

	ParamError := err.(RadioParamError)
	e.Add([]string{ParamError.Field}, "ParamError", ParamError.Msg)
	return e
}

func checkRadioParam(r *http.Request, param string) error {
	p := r.FormValue(param)
	if p != "" {
		if p != "on" {
			return RadioParamError{param, "invalid param."}
		}
		return nil
	}
	return nil
}

func (f NewPairForm) FormatTo(to interface{}) error {
	config := to.(*settings.SyncConfig)
	config.SyncDelete = radio[f.SyncDel]
	config.SyncRename = radio[f.SyncRename]
	config.CoverSameName = radio[f.CoverSameName]
	return nil
}

var radio = make(map[string]bool, 3)

func init() {
	radio["on"] = true
	radio[""] = false
}
