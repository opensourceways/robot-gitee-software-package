package watchingimpl

type Config struct {
	Org       string    `json:"org" required:"true"`
	Frequency Frequency `json:"frequency"`
}

type Frequency struct {
	MaxTimes int `json:"max_times"`
	Interval int `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.Frequency.MaxTimes <= 0 {
		cfg.Frequency.MaxTimes = 3
	}
	if cfg.Frequency.Interval <= 0 {
		cfg.Frequency.Interval = 10
	}
}
