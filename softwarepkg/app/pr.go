package app

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/email"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/watching"
)

type PullRequestService interface {
	HandleCI(cmd *CmdToHandleCI) error
}

type pullRequestService struct {
	repo     repository.PullRequest
	producer message.SoftwarePkgMessage
	email    email.Email
	watch    watching.Watching
	log      *logrus.Entry
}

func (s *pullRequestService) HandleCI(cmd *CmdToHandleCI) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if !cmd.isSuccess() {
		if err = s.email.Send(pr.Link); err != nil {
			return err
		}
	}

	e := domain.NewPRCIFinishedEvent(&pr, cmd.FailedReason)
	return s.producer.NotifyCIResult(&e)
}

func (s *pullRequestService) watchCreateRepo(pr domain.PullRequest) {
	v := domain.ToSoftwarePkgRepo(&pr)

	if err := s.watch.Apply(v); err != nil {
		s.log.Error(err)

		return
	}

	e := domain.NewCreateRepoEvent(v)
	if err := s.producer.NotifyCreateRepoResult(&e); err != nil {
		s.log.WithError(err).Error("notify create repo event failed")

		return
	}

	if err := s.repo.Remove(&pr); err != nil {
		s.log.WithError(err).Error("remove pr storage failed")
	}
}

func (s *pullRequestService) InitWatchCreateRepo() error {
	prs, err := s.repo.FindAll()
	if err != nil {
		return err
	}

	for _, pr := range prs {
		v := pr
		if !v.IsMerged() {
			continue
		}

		go s.watchCreateRepo(v)
	}

	return nil
}
