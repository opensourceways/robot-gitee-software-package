package client

import sdk "github.com/opensourceways/go-gitee/gitee"

type IClient interface {
	GetBot() (sdk.User, error)
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
}
