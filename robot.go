package main

import (
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/framework"
	"github.com/opensourceways/server-common-lib/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

// TODO: set botName
const botName = "software-package"

type iClient interface {
}

func newRobot(cli iClient, prService *app.PullRequestService) *robot {
	return &robot{
		cli:       cli,
		prService: prService,
	}
}

type robot struct {
	cli       iClient
	prService *app.PullRequestService
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config, org, repo string) (*botConfig, error) {
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, fmt.Errorf("can't convert to configuration")
	}

	if bc := c.configFor(org, repo); bc != nil {
		return bc, nil
	}

	return nil, fmt.Errorf("no config for this repo:%s/%s", org, repo)
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	f.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	org, repo := e.GetOrgRepo()
	cfg, err := bot.getConfig(c, org, repo)
	if err != nil {
		return err
	}

	prState := e.GetPullRequest().GetState()

	if prState == sdk.StatusOpen &&
		sdk.GetPullRequestAction(e) == sdk.PRActionUpdatedLabel {

		return bot.handleCILabel(e, cfg)
	}

	return nil
}

func (bot *robot) handleCILabel(e *sdk.PullRequestEvent, cfg *botConfig) error {
	labels := e.PullRequest.LabelsToSet()

	cmd := app.CmdToHandleCI{
		PRNum:        int(e.Number),
		FailedReason: "",
	}

	if labels.Has(cfg.CILabel.Success) {
		if err := bot.prService.HandleCI(&cmd); err != nil {
			return err
		}
	}

	if labels.Has(cfg.CILabel.Fail) {
		cmd.FailedReason = "ci check failed"
		if err := bot.prService.HandleCI(&cmd); err != nil {
			return err
		}

		return bot.sendEmail(cfg, e.GetURL())
	}

	return nil
}

func (bot *robot) sendEmail(cfg *botConfig, url string) error {
	d := gomail.NewDialer(
		cfg.EmailServer.Host,
		cfg.EmailServer.Port,
		cfg.EmailServer.From,
		cfg.EmailServer.AuthCode)

	subject := "the CI of PR in openeuler/community is failed"
	content := fmt.Sprintf("the pr url: %s", url)

	message := gomail.NewMessage()
	message.SetHeader("From", cfg.EmailServer.From)
	message.SetHeader("To", cfg.MaintainerEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", content)

	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
