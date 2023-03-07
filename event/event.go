package event

import (
	"encoding/json"
	"fmt"

	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/client"
)

var robotLogin string

type Event struct {
	cli   client.IClient
	cfg   *Config
	log   *logrus.Entry
	group string
}

func NewEvent(cfg *Config, cli client.IClient, group string) *Event {
	return &Event{
		cli:   cli,
		cfg:   cfg,
		group: group,
	}
}

func (e *Event) Subscribe() (subscribers map[string]mq.Subscriber, err error) {
	subscribers = make(map[string]mq.Subscriber)

	s, err := kafka.Subscribe(e.cfg.Topics.NewPkg, e.group, e.newPkgHandle)
	if err != nil {
		return
	}
	subscribers[s.Topic()] = s

	return
}

func (e *Event) newPkgHandle(event mq.Event) error {
	e.log = logrus.WithFields(
		logrus.Fields{
			"msg": event.Message(),
		},
	)

	e.createPR(event.Message())

	return nil
}

func (e *Event) createPR(msg *mq.Message) {
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

func (e *Event) createPRWithApi(p CreatePRParam) error {
	robotName, err := e.getRobotLogin()
	if err != nil {
		return err
	}

	head := fmt.Sprintf("%s:%s", robotName, branchName(e.cfg.PR.BranchName, p.PkgName))
	pr, err := e.cli.CreatePullRequest(
		e.cfg.PR.Org, e.cfg.PR.Repo, prName(e.cfg.PR.PRName, p.PkgName),
		p.ReasonToImportPkg, head, "master", true,
	)
	if err != nil {
		return err
	}

	e.log.Infof("pr number is %d", pr.Number)

	return nil
}

func (e *Event) getRobotLogin() (string, error) {
	if robotLogin == "" {
		v, err := e.cli.GetBot()
		if err != nil {
			return "", err
		}

		robotLogin = v.Login
	}

	return robotLogin, nil
}

func branchName(branchName, pkgName string) string {
	return fmt.Sprintf(branchName, pkgName)
}

func prName(prName, pkgName string) string {
	return fmt.Sprintf(prName, pkgName)
}
