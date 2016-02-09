package settings

import (
	// "github.com/Felamande/filesync/server/modules/utils"
	"crypto/md5"
	"encoding/hex"
	"io"
    "errors"
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
	Left      string     `json:"left"`
	Right     string     `json:"right"`
	Config    SyncConfig `json:"config"`
	IgnoreExt []string   `yaml:"ignore_ext"`
}

type cfgMgr struct {
	cfg    *SavedConfig
	hashMap map[string]bool
}

func newMgr()*cfgMgr{
    return &cfgMgr{
        hashMap:make(map[string]bool),
        cfg
    }
}

func (m *cfgMgr) Cfg() *SavedConfig {
	if m.cfg != nil {
		return m.cfg
	}
	return new(SavedConfig)
}

func (m *cfgMgr) init() {
    s := readConfig(getAbs(settingStruct.Filesync.CfgFile))
	for _, val := range s.Pairs {
		m.Add(val)
	}

}

func (m *cfgMgr) Add(p *SyncPairConfig) error {
	hash := mD5(p.Left, p.Right)
	if m.hashMap[hash]{
        return errors.New("duplicated pair")
    }
	m.cfg.Pairs = append(m.cfg.Pairs, p)
	return nil
}

func (m *cfgMgr) Save() error {
	
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
