package main

import (
	"encoding/json"
	"fmt"

	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
)

const (
	org = "openeuler"

	repoName = "community"

	eventCreatePR = "create_PR"

	branchNameFormat = "software_pkg_%s"

	prNameFormat = branchNameFormat + "新增软件包申请"
)

var robotLogin string

type event struct {
	cli iClient
	cfg *configuration
	log *logrus.Entry
}

func newEvent(cfg *configuration, cli iClient) *event {
	return &event{
		cli: cli,
		cfg: cfg,
	}
}

func (e *event) process() (mq.Subscriber, error) {
	if err := e.init(); err != nil {
		return nil, err
	}

	return kafka.Subscribe(e.cfg.Topic, botName, e.handle)
}

func (e *event) init() error {
	v, err := e.cli.GetBot()
	if err != nil {
		return err
	}

	robotLogin = v.Login

	return nil
}

func (e *event) handle(event mq.Event) error {
	eventType := event.Message().Header["event_type"]

	e.log = logrus.WithFields(
		logrus.Fields{
			"event_type":  eventType,
			"msg_content": event.Message(),
		},
	)

	switch eventType {
	case eventCreatePR:
		e.createPR(event.Message())
	default:

	}

	return nil
}

func (e *event) createPR(msg *mq.Message) {
	var c CreatePRParam
	if err := json.Unmarshal(msg.Body, &c); err != nil {
		e.log.WithError(err).Error("unmarshal")
		return
	}

	if err := c.initRepo(e.cfg); err != nil {
		e.log.WithError(err).Error("init repo")
		return
	}

	if err := c.newBranch(e.cfg); err != nil {
		e.log.WithError(err).Error("new branch")
		return
	}

	if err := c.modifyFiles(e.cfg); err != nil {
		e.log.WithError(err).Error("modify files")
		return
	}

	if err := c.commit(e.cfg); err != nil {
		e.log.WithError(err).Error("commit")
		return
	}

	if err := e.createPRWithApi(c); err != nil {
		e.log.WithError(err).Error("create with api")
		return
	}
}

func (e *event) createPRWithApi(p CreatePRParam) error {
	head := fmt.Sprintf("%s:%s", robotLogin, branchName(p.PackageName))
	pr, err := e.cli.CreatePullRequest(org, repoName, prName(p.PackageName), p.Purpose, head, "master", true)
	if err != nil {
		return err
	}

	logrus.Infof("pr number is %d", pr.Number)

	return nil
}

func branchName(pkgName string) string {
	return fmt.Sprintf(branchNameFormat, pkgName)
}

func prName(pkgName string) string {
	return fmt.Sprintf(prNameFormat, pkgName)
}
