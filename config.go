package main

import (
	"github.com/opensourceways/server-common-lib/config"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-software-package/message"
)

type configuration struct {
	ConfigItems []botConfig    `json:"config_items,omitempty"`
	Event       message.Config `json:"event_config"`
}

func LoadConfig(path string) (*configuration, error) {
	cfg := new(configuration)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *configuration) configFor(org, repo string) *botConfig {
	if c == nil {
		return nil
	}

	items := c.ConfigItems
	v := make([]config.IRepoFilter, len(items))
	for i := range items {
		v[i] = &items[i]
	}

	if i := config.Find(org, repo, v); i >= 0 {
		return &items[i]
	}

	return nil
}

func (c *configuration) Validate() error {
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

func (c *configuration) SetDefault() {
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
