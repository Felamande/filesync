package pairs

import(
    // "github.com/lunny/tango"
    // "errors"
    "github.com/Felamande/filesync/syncer"
    "github.com/Felamande/filesync/server/routers/base"
)

type NewPairRouter struct{
    base.BaseJSONRouter
    
}

func (r *NewPairRouter)Get()interface{}{
    return syncer.Default().PairMap
}
