package models

import "fmt"

type RadioParamError struct{
    Field string
    Msg   string
}

func(e RadioParamError)Error()string{
    return fmt.Sprintf("in field %s: %s",e.Field,e.Msg)
}