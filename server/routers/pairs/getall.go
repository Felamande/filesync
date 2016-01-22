package pairs

import(
    // "github.com/lunny/tango"
    // "errors"
    "github.com/Felamande/filesync/syncer"
    "github.com/Felamande/filesync/server/routers/base"
)

type GetAllRouter struct{
    base.BaseJSONRouter
    
}

func (r *GetAllRouter)Get()interface{}{
    return syncer.Default().PairMap
}
