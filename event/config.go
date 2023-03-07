package event

import (
	"errors"
)

type Config struct {
	KafkaAddress  string        `json:"kafka_address"   required:"true"`
	Topics        Topics        `json:"topics"`
	Robot         RobotConfig   `json:"robot"`
	PkgRepoBranch PkgRepoBranch `json:"pkg_repo_branch"`
}

func (c *Config) Validate() error {
	if c.KafkaAddress == "" {
		return errors.New("missing kafka_address")
	}

	if c.Topics.NewPkg == "" {
		return errors.New("missing new pkg topic")
	}

	if c.Topics.CIPassed == "" {
		return errors.New("missing ci passed topic")
	}

	return nil
}

func (c *Config) SetDefault() {
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

type PkgRepoBranch struct {
	Name        string `json:"name"`
	ProtectType string `json:"protect_type"`
	PublicType  string `json:"public_type"`
}

type RobotConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Topics struct {
	NewPkg   string `json:"new_pkg"`
	CIPassed string `json:"ci_passed"`
}
