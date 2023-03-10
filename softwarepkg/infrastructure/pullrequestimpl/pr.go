package pullrequestimpl

import (
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"gopkg.in/gomail.v2"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type pullRequestImpl struct {
	cli        iClient
	cfg        Config
	pkg        *domain.SoftwarePkg
	robotLogin string
}

type iClient interface {
	GetBot() (sdk.User, error)
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
}

func (impl *pullRequestImpl) Create(pkg *domain.SoftwarePkg) (pr domain.PullRequest, err error) {
	impl.pkg = pkg

	if err = impl.initRepo(); err != nil {
		return
	}

	if err = impl.newBranch(); err != nil {
		return
	}

	if err = impl.modifyFiles(); err != nil {
		return
	}

	if err = impl.commit(); err != nil {
		return
	}

	return impl.submit()
}

func (impl *pullRequestImpl) Merge(*domain.PullRequest) error {
	return nil
}

func (impl *pullRequestImpl) Close(*domain.PullRequest) error {
	return nil
}

func (impl *pullRequestImpl) SendEmail(url string) error {
	d := gomail.NewDialer(
		impl.cfg.EmailServer.Host,
		impl.cfg.EmailServer.Port,
		impl.cfg.EmailServer.From,
		impl.cfg.EmailServer.AuthCode,
	)

	subject := "the CI of PR in openeuler/community is failed"
	content := fmt.Sprintf("the pr url: %s", url)

	message := gomail.NewMessage()
	message.SetHeader("From", impl.cfg.EmailServer.From)
	message.SetHeader("To", impl.cfg.MaintainerEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", content)

	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
