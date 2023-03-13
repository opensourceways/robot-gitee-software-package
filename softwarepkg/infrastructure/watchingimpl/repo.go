package watchingimpl

import (
	"errors"
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func NewWatchingImpl(cfg Config, cli iClient) *WatchingImpl {
	return &WatchingImpl{
		cfg: cfg,
		cli: cli,
	}
}

type WatchingImpl struct {
	cfg Config
	cli iClient
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

func (impl WatchingImpl) Apply(pkg *domain.SoftwarePkgRepo) error {
	for i := 0; i < impl.cfg.Frequency.MaxTimes; i++ {
		time.Sleep(time.Duration(impl.cfg.Frequency.Interval) * time.Second)

		repo, err := impl.cli.GetRepo(impl.cfg.Org, pkg.Pkg.Name)
		if err == nil {
			pkg.RepoURL = repo.Url

			return nil
		}
	}

	return errors.New("repo hasn't been created")
}
