package syncer

type GlobalConfig struct {
}

type SavedConfig struct {
	Pairs []SyncPairConfig `json:"pairs"`
	Port  string           `json:"port"`
}

type SyncPairConfig struct {
	Left   string     `json:"left"`
	Right  string     `json:"right`
	Config SyncConfig `json:"config"`
}
