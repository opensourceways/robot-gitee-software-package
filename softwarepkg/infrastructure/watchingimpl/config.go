package watchingimpl

type Config struct {
	Org      string `json:"org" required:"true"`
	Interval int    `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
}
