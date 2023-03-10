package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type MessageService interface {
	CreatePR(cmd *CmdToCreatePR) error
}

type PullRequestService struct {
	repo     repository.PullRequest
	prCli    pullrequest.PullRequest
	producer message.SoftwarePkgMessage
}

func (s *PullRequestService) CreatePR(cmd *CmdToCreatePR) error {
	pr, err := s.prCli.Create(cmd)
	if err != nil {
		return err
	}

	return s.repo.Add(&pr)
}

func (s *PullRequestService) HandleCI(cmd *CmdToHandleCI) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	e := domain.NewPRCIFinishedEvent(&pr, cmd.FailedReason)
	return s.producer.NotifyCIResult(&e)
}
