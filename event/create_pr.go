package event

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
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

type CreatePRParam domain.SoftwarePkgAppliedEvent

func (c CreatePRParam) modifyFiles(cfg *Config) error {
	if err := c.appendToSigInfo(); err != nil {
		return err
	}

	return c.newCreateRepoYaml(cfg)
}

func (c CreatePRParam) appendToSigInfo() error {
	appendContent := fmt.Sprintf(appendToSigInfo, c.PkgName, c.ImporterEmail, c.Importer)
	fileName := fmt.Sprintf("community/sig/%s/sig-info.yaml", c.ImportingPkgSig)

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

func (c CreatePRParam) newCreateRepoYaml(cfg *Config) error {
	subDirName := strings.ToLower(c.PkgName[:1])
	fileName := fmt.Sprintf("community/sig/%s/src-openeuler/%s/%s.yaml",
		c.ImportingPkgSig, subDirName, c.PkgName,
	)

	content := fmt.Sprintf(createRepoConfigFile,
		c.PkgName, c.PkgDesc, c.SourceCodeURL,
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

func (c CreatePRParam) initRepo(cfg *Config) error {
	if s, err := os.Stat(repoName); err == nil && s.IsDir() {
		return nil
	}

	return c.execScript(cfg, cmdInit)
}

func (c CreatePRParam) newBranch(cfg *Config) error {
	return c.execScript(cfg, cmdNewBranch)
}

func (c CreatePRParam) commit(cfg *Config) error {
	return c.execScript(cfg, cmdCommit)
}

func (c CreatePRParam) execScript(cfg *Config, cmdType CmdType) error {
	cmd := exec.Command(repoHandleScript, string(cmdType), cfg.Robot.Username,
		cfg.Robot.Password, cfg.Robot.Email, branchName(c.PkgName))

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}
