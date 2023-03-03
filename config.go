package main

import (
	"errors"
)

type configuration struct {
	Topic        string       `json:"topic"           required:"true"`
	KafkaAddress string       `json:"kafka_address"   required:"true"`
	Branch       BranchConfig `json:"branch"`
	Robot        RobotConfig  `json:"robot"`
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
	if c.Branch.Name == "" {
		c.Branch.Name = "master"
	}

	if c.Branch.ProtectType == "" {
		c.Branch.ProtectType = "protected"
	}

	if c.Branch.PublicType == "" {
		c.Branch.PublicType = "public"
	}
}

type BranchConfig struct {
	Name        string `json:"name"`
	ProtectType string `json:"protect_type"`
	PublicType  string `json:"public_type"`
}

type RobotConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
