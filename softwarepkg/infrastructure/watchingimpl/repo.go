package watchingimpl

import (
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func NewWatchingImpl(cfg Config, cli iClient) *WatchingImpl {
	return &WatchingImpl{
		cfg: cfg,
		cli: cli,
	}
}

type WatchingImpl struct {
	cfg       Config
	cli       iClient
	log       *logrus.Entry
	repo      repository.PullRequest
	prService app.PullRequestService
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

func (impl *WatchingImpl) Run() {
	for {
		prs, err := impl.repo.FindAll()
		if err != nil {
			impl.log.WithError(err).Error("find all storage pr failed")
		}

		for _, pr := range prs {
			if !pr.IsMerged() {
				continue
			}

			v, err := impl.cli.GetRepo(impl.cfg.Org, pr.Pkg.Name)
			if err != nil {
				continue
			}

			impl.prService.HandleRepoCreated(&pr, v.Url)
		}

		time.Sleep(time.Second * time.Duration(impl.cfg.Interval))
	}
}
