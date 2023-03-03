package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
)

const (
	org = "openeuler"

	repoName = "community"

	eventCreatePR = "create_PR"

	repoHandleScript = "./repo.sh"

	branchNameFormat = "software_pkg_%s"

	prNameFormat = branchNameFormat + "新增软件包申请"
)

var appendToSigInfo = `
- repo:
  - src-openeuler/%s
  committers:
  - email: %s
    name: %s
`

var createRepoConfigFile = `
name: %s
description: %s
upstream: %s
branches:
- name: %s
type: %s
type: %s
`

type event struct {
	cli       iClient
	cfg       *configuration
	eventType string
}

func newEvent(cfg *configuration, cli iClient) *event {
	return &event{
		cli: cli,
		cfg: cfg,
	}
}

func (e *event) process() (mq.Subscriber, error) {
	return kafka.Subscribe(e.cfg.Topic, botName, e.handle)
}

func (e *event) handle(event mq.Event) error {
	e.eventType = event.Message().Header["event_type"]
	switch e.eventType {
	case eventCreatePR:
		e.createPR(event.Message())
	default:

	}

	return nil
}

func (e *event) createPR(msg *mq.Message) {
	l := logrus.WithFields(
		logrus.Fields{
			"event-type": e.eventType,
		},
	)

	var param CreatePRParam
	if err := json.Unmarshal(msg.Body, &param); err != nil {
		l.WithError(err).Error("unmarshal")
		return
	}

	if err := param.initRepo(e.cfg); err != nil {
		l.WithError(err).Error("init repo")
		return
	}

	if err := param.newBranch(e.cfg); err != nil {
		l.WithError(err).Error("new branch")
		return
	}

	if err := param.modifyFiles(e.cfg); err != nil {
		l.WithError(err).Error("modify files")
		return
	}

	if err := param.commit(e.cfg); err != nil {
		l.WithError(err).Error("commit")
		return
	}

	if err := e.createPRWithApi(param); err != nil {
		l.WithError(err).Error("create with api")
		return
	}
}

func (e *event) createPRWithApi(p CreatePRParam) error {
	v, err := e.cli.GetBot()
	if err != nil {
		return err
	}

	head := fmt.Sprintf("%s:%s", v.Login, p.branchName())
	pr, err := e.cli.CreatePullRequest(org, repoName, p.prName(), p.Purpose, head, "master", true)
	if err != nil {
		return err
	}

	logrus.Infof("pr number is %d", pr.Number)

	return nil
}

type CreatePRParam struct {
	User        User   `json:"user"`
	PackageName string `json:"package_name"`
	Description string `json:"description"`
	Purpose     string `json:"purpose"`
	Upstream    string `json:"upstream"`
	SIG         string `json:"sig"`
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (p CreatePRParam) modifyFiles(cfg *configuration) error {
	if err := p.appendToSigInfo(); err != nil {
		return err
	}

	return p.newCreateRepoYaml(cfg)
}

func (p CreatePRParam) appendToSigInfo() error {
	appendContent := fmt.Sprintf(appendToSigInfo, p.PackageName, p.User.Email, p.User.Name)
	fileName := fmt.Sprintf("community/sig/%s/sig-info.yaml", p.SIG)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	if _, err = write.WriteString(appendContent); err != nil {
		return err
	}

	if err = write.Flush(); err != nil {
		return err
	}

	return nil
}

func (p CreatePRParam) newCreateRepoYaml(cfg *configuration) error {
	subDirName := strings.ToLower(p.PackageName[:1])
	fileName := fmt.Sprintf("community/sig/%s/src-openeuler/%s/%s.yaml",
		p.SIG, subDirName, p.PackageName,
	)

	content := fmt.Sprintf(createRepoConfigFile,
		p.PackageName, p.Description, p.Upstream,
		cfg.Branch.Name, cfg.Branch.ProtectType, cfg.Branch.PublicType,
	)

	return os.WriteFile(fileName, []byte(content), 0644)
}

type CmdType string

var (
	cmdInit      = CmdType("init")
	cmdNewBranch = CmdType("new")
	cmdCommit    = CmdType("commit")
)

func (p CreatePRParam) initRepo(cfg *configuration) error {
	if _, err := os.Stat(repoName); err == nil {
		return nil
	}

	return p.execScript(cfg, cmdInit)
}

func (p CreatePRParam) newBranch(cfg *configuration) error {
	return p.execScript(cfg, cmdNewBranch)
}

func (p CreatePRParam) commit(cfg *configuration) error {
	return p.execScript(cfg, cmdCommit)
}

func (p CreatePRParam) execScript(cfg *configuration, cmdType CmdType) error {
	cmd := exec.Command(repoHandleScript, string(cmdType), cfg.Robot.Username,
		cfg.Robot.Password, p.branchName())

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (p CreatePRParam) branchName() string {
	return fmt.Sprintf(branchNameFormat, p.PackageName)
}

func (p CreatePRParam) prName() string {
	return fmt.Sprintf(prNameFormat, p.PackageName)
}
