package watchingimpl

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
)

func NewWatchingImpl(cfg Config, cli iClient) *WatchingImpl {
	return &WatchingImpl{
		cfg:      cfg,
		cli:      cli,
		repoChan: make(chan *domain.SoftwarePkgRepo, 100),
	}
}

type WatchingImpl struct {
	cfg      Config
	cli      iClient
	repoChan chan *domain.SoftwarePkgRepo
	producer message.SoftwarePkgMessage
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

func (impl *WatchingImpl) Apply(pkg *domain.SoftwarePkgRepo) error {
	// TODO some validate

	impl.repoChan <- pkg

	return nil
}

func (impl *WatchingImpl) Run() {
	for {
		select {
		case repo := <-impl.repoChan:
			v, err := impl.cli.GetRepo(impl.cfg.Org, repo.Pkg.Name)
			if err != nil {
				impl.repoChan <- repo
				continue
			}

			repo.RepoURL = v.Url
			impl.producer.NotifyRepoCreatedResult(repo)
		}
	}
}
