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

	branchNameFormat = "software_pkg_%s"

	prNameFormat = branchNameFormat + ",新增软件包申请"
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

func (e *event) subscribe() (subscribers map[string]mq.Subscriber, err error) {
	subscribers = make(map[string]mq.Subscriber)

	s, err := kafka.Subscribe(e.cfg.Topics.NewPkg, botName, e.newPkgHandle)
	if err != nil {
		return
	}
	subscribers[s.Topic()] = s

	return
}

func (e *event) newPkgHandle(event mq.Event) error {
	e.log = logrus.WithFields(
		logrus.Fields{
			"msg": event.Message(),
		},
	)

	e.createPR(event.Message())

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
	robotName, err := e.getRobotLogin()
	if err != nil {
		return err
	}

	head := fmt.Sprintf("%s:%s", robotName, branchName(p.PkgName))
	pr, err := e.cli.CreatePullRequest(
		org, repoName, prName(p.PkgName),
		p.ReasonToImportPkg, head, "master", true,
	)
	if err != nil {
		return err
	}

	logrus.Infof("pr number is %d", pr.Number)

	return nil
}

func (e *event) getRobotLogin() (string, error) {
	if robotLogin == "" {
		v, err := e.cli.GetBot()
		if err != nil {
			return "", err
		}

		robotLogin = v.Login
	}

	return robotLogin, nil
}

func branchName(pkgName string) string {
	return fmt.Sprintf(branchNameFormat, pkgName)
}

func prName(pkgName string) string {
	return fmt.Sprintf(prNameFormat, pkgName)
}
