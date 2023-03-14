package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/watchingimpl"
)

type MessageService interface {
	CreatePR(cmd *CmdToCreatePR) error
	MergePR(cmd *CmdToMergePR) error
	ClosePR(cmd *CmdToClosePR) error
}

type messageService struct {
	repo  repository.PullRequest
	prCli pullrequest.PullRequest
	watch watchingimpl.WatchingImpl
}

func (s *messageService) CreatePR(cmd *CmdToCreatePR) error {
	pr, err := s.prCli.Create(cmd)
	if err != nil {
		return err
	}

	return s.repo.Add(&pr)
}

func (s *messageService) MergePR(cmd *CmdToMergePR) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if err = s.prCli.Merge(&pr); err != nil {
		return err
	}

	v := domain.ToSoftwarePkgRepo(&pr)
	return s.watch.Apply(v)
}

func (s *messageService) ClosePR(cmd *CmdToClosePR) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	return s.prCli.Close(&pr)
}
