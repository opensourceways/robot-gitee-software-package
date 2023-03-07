package config

import (
	"github.com/opensourceways/server-common-lib/config"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-software-package/event"
)

type Config struct {
	ConfigItems []botConfig  `json:"config_items,omitempty"`
	Event       event.Config `json:"event_config"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}

	items := c.ConfigItems
	for i := range items {
		if err := items[i].validate(); err != nil {
			return err
		}
	}

	if err := c.Event.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *Config) SetDefault() {
	if c == nil {
		return
	}

	Items := c.ConfigItems
	for i := range Items {
		Items[i].setDefault()
	}

	c.Event.SetDefault()
}

type botConfig struct {
	config.RepoFilter
}

func (c *botConfig) setDefault() {
}

func (c *botConfig) validate() error {
	return c.RepoFilter.Validate()
}
