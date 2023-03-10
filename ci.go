package main

import (
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"gopkg.in/gomail.v2"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

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
