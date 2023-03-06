package main

import (
	"errors"
)

type configuration struct {
	Topic         string              `json:"topic"           required:"true"`
	KafkaAddress  string              `json:"kafka_address"   required:"true"`
	Robot         RobotConfig         `json:"robot"			  required:"true"`
	PkgRepoBranch PkgRepoBranchConfig `json:"pkg_repo_branch"`
}

func (c *configuration) Validate() error {
	if c.Topic == "" {
		return errors.New("missing topic")
	}

	if c.KafkaAddress == "" {
		return errors.New("missing kafka_address")
	}

	return nil
}

func (c *configuration) SetDefault() {
	if c.PkgRepoBranch.Name == "" {
		c.PkgRepoBranch.Name = "master"
	}

	if c.PkgRepoBranch.ProtectType == "" {
		c.PkgRepoBranch.ProtectType = "protected"
	}

	if c.PkgRepoBranch.PublicType == "" {
		c.PkgRepoBranch.PublicType = "public"
	}
}

type PkgRepoBranchConfig struct {
	Name        string `json:"name"`
	ProtectType string `json:"protect_type"`
	PublicType  string `json:"public_type"`
}

type RobotConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
