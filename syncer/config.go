package syncer

type GlobalConfig struct {
}

type SavedConfig struct {
	Pairs   []SyncPairConfig `json:"pairs"`
	Port    int              `json:"port"`
	LogPath string           `json:"log_path"`
	Ignore  []string         `json:"ignore_rules"`
}

type SyncPairConfig struct {
	Left   string     `json:"left"`
	Right  string     `json:"right`
	Config SyncConfig `json:"config"`
}
