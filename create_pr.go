package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const repoHandleScript = "./repo.sh"

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

func (c CreatePRParam) modifyFiles(cfg *configuration) error {
	if err := c.appendToSigInfo(); err != nil {
		return err
	}

	return c.newCreateRepoYaml(cfg)
}

func (c CreatePRParam) appendToSigInfo() error {
	appendContent := fmt.Sprintf(appendToSigInfo, c.PackageName, c.User.Email, c.User.Name)
	fileName := fmt.Sprintf("community/sig/%s/sig-info.yaml", c.SIG)

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

func (c CreatePRParam) newCreateRepoYaml(cfg *configuration) error {
	subDirName := strings.ToLower(c.PackageName[:1])
	fileName := fmt.Sprintf("community/sig/%s/src-openeuler/%s/%s.yaml",
		c.SIG, subDirName, c.PackageName,
	)

	content := fmt.Sprintf(createRepoConfigFile,
		c.PackageName, c.Description, c.Upstream,
		cfg.PkgRepoBranch.Name,
		cfg.PkgRepoBranch.ProtectType,
		cfg.PkgRepoBranch.PublicType,
	)

	return os.WriteFile(fileName, []byte(content), 0644)
}

type CmdType string

var (
	cmdInit      = CmdType("init")
	cmdNewBranch = CmdType("new")
	cmdCommit    = CmdType("commit")
)

func (c CreatePRParam) initRepo(cfg *configuration) error {
	if s, err := os.Stat(repoName); err == nil && s.IsDir() {
		return nil
	}

	return c.execScript(cfg, cmdInit)
}

func (c CreatePRParam) newBranch(cfg *configuration) error {
	return c.execScript(cfg, cmdNewBranch)
}

func (c CreatePRParam) commit(cfg *configuration) error {
	return c.execScript(cfg, cmdCommit)
}

func (c CreatePRParam) execScript(cfg *configuration, cmdType CmdType) error {
	cmd := exec.Command(repoHandleScript, string(cmdType), cfg.Robot.Username,
		cfg.Robot.Password, cfg.Robot.Email, branchName(c.PackageName))

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}
