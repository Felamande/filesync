package settings

import (
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
    hash      string     
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

func (m *cfgMgr) init(s *SavedConfig) {
	m.schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"pairs": &memdb.TableSchema{
				Name: "pairs",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Left"},
					},
					"left": &memdb.IndexSchema{
						Name:    "left",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Left"},
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
}

func (m *cfgMgr) Add(p *SyncPairConfig) error {
	txn := m.db.Txn(true)
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

	r, err := txn.Get("pairs", "left")
	if err != nil {
		return err
	}
	for rr := r.Next(); rr != nil; rr = r.Next() {
		if s, ok := rr.(*SyncPairConfig); ok {
			continue
		} else {
			m.cfg.Pairs = append(m.cfg.Pairs, s)
		}
	}
    return nil

}
