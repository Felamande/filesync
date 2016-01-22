package syncer

type GlobalConfig struct {
}

type SavedConfig struct {
	Pairs     []SyncPairConfig `json:"pairs"`
	LogPath   string           `json:"log_path"`
	IgnoreExt []string         `yaml:"ignore_ext"`
}

type SyncPairConfig struct {
	Left   string     `json:"left"`
	Right  string     `json:"right`
	Config SyncConfig `json:"config"`
}
