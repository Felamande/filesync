package settings

import (
	// "github.com/Felamande/filesync/server/modules/utils"
	"encoding/hex"
	"io"
    "crypto/md5"

	memdb "github.com/hashicorp/go-memdb"
)

type SavedConfig struct {
	Pairs []*SyncPairConfig `json:"pairs"`
}

type SyncConfig struct {
	CoverSameName bool `json:"cover_same_name"`
	SyncDelete    bool `json:"sync_delete"`
	SyncRename    bool `json:"sync_rename"`
}

type SyncPairConfig struct {
	Hash      string     `json:"-" yaml:"-"`
	Left      string     `json:"left"`
	Right     string     `json:"right"`
	IgnoreExt []string   `yaml:"ignore_ext"`
	Config    SyncConfig `json:"config"`
}

type cfgMgr struct {
	cfg    *SavedConfig
	schema *memdb.DBSchema
	db     *memdb.MemDB
}

func (m *cfgMgr) Init(s *SavedConfig) {
	m.schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"pairs": &memdb.TableSchema{
				Name: "pairs",
				Indexes: map[string]*memdb.IndexSchema{
					"hash": &memdb.IndexSchema{
						Name:    "hash",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Hash"},
					},
				},
			},
		},
	}
	m.cfg = s
	db, err := memdb.NewMemDB(m.schema)
	if err != nil {
		panic(err)
	}
	m.db = db
	for _, val := range s.Pairs {
		m.Add(val)
	}

}

func (m *cfgMgr) Add(p *SyncPairConfig) error {
	txn := m.db.Txn(true)
	p.Hash = mD5(p.Left, p.Right)
	err := txn.Insert("pairs", p)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (m *cfgMgr) Save() error {
	txn := m.db.Txn(false)
	defer txn.Abort()

	r, err := txn.Get("pairs", "hash")
	if err != nil {
		return err
	}
    saved:=new(SavedConfig)
	for rr := r.Next(); rr != nil; rr = r.Next() {
		if s, ok := rr.(*SyncPairConfig); !ok {
			continue
		} else {
			saved.Pairs = append(saved.Pairs, s)
		}
	}
    
	return nil

}

func mD5(source ...interface{}) string {
	ctx := md5.New()
	for _, s := range source {
		switch ss := s.(type) {
		case io.Reader:
			io.Copy(ctx, ss)
		case string:
			ctx.Write([]byte(ss))
		case []byte:
			ctx.Write(ss)

		}
	}

	return hex.EncodeToString(ctx.Sum(nil))
}
