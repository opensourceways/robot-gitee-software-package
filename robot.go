package main

import (
	"errors"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/server-common-lib/config"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-lib/framework"

	"github.com/opensourceways/robot-gitee-software-package/client"
	config2 "github.com/opensourceways/robot-gitee-software-package/config"
)

const botName = "software-package"

func newRobot(cli client.IClient) *robot {
	return &robot{cli: cli}
}

type robot struct {
	cli client.IClient
}

func (bot *robot) NewConfig() config.Config {
	return &config2.Config{}
}

func (bot *robot) getConfig(cfg config.Config) (*config2.Config, error) {
	if c, ok := cfg.(*config2.Config); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	f.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand PR event, delete this function.
	return nil
}
