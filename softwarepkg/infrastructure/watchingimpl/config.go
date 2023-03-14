package watchingimpl

type Config struct {
	Org      string `json:"org"`
	Interval int    `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.Org == "" {
		cfg.Org = "src-openeuler"
	}

	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
}
