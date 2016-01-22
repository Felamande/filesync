package time

import(
    "time"
    "github.com/lunny/tango"
)

type TimeHandler struct{
        
}

func(h *TimeHandler) Handle(ctx *tango.Context){
    t1:=time.Now()
    ctx.Next()
    ctx.Logger.Infof("Completed %v %v %v in %v for %v",ctx.Status(),ctx.Req().Method,ctx.Req().URL.Path,time.Since(t1),ctx.Req().RemoteAddr)
}