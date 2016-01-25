package settings

import (
	memdb "github.com/hashicorp/go-memdb"
)


var pairSchema = &memdb.DBSchema{
    Tables:map[string]*memdb.TableSchema{
        "pairs":&memdb.TableSchema{
            Name:"pairs",
            Indexes:map[string]*memdb.IndexSchema{
                "id":&memdb.IndexSchema{
                    Name:"id",
                    Unique:true,
                    Indexer:&memdb.StringFieldIndex{Field:"Left"},
                },
                "left":&memdb.IndexSchema{
                    Name:"left",
                    Unique:false,
                    Indexer:&memdb.StringFieldIndex{Field:"Left"},
                },
            },  
        },
    },
}




func init(){
    
}