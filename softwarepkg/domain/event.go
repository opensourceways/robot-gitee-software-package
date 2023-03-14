package domain

import "encoding/json"

type prCIFinishedEvent struct {
	PkgId        string `json:"pkg_id"`
	PkgName      string `json:"pkg_name"`
	PRLink       string `json:"pr_link"`
	FailedReason string `json:"failed_reason"`
	Success      bool   `json:"success"`
}

func (e *prCIFinishedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRCIFinishedEvent(pr *PullRequest, failedReason string) prCIFinishedEvent {
	return prCIFinishedEvent{
		PkgId:        pr.Pkg.Id,
		PkgName:      pr.Pkg.Name,
		PRLink:       pr.Link,
		FailedReason: failedReason,
		Success:      failedReason == "",
	}
}

type repoCreatedEvent struct {
	PkgId   string `json:"pkg_id"`
	PkgName string `json:"pkg_name"`
	RepoURL string `json:"repo_url"`
}

func (e *repoCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewRepoCreatedEvent(pr *PullRequest, url string) repoCreatedEvent {
	return repoCreatedEvent{
		PkgId:   pr.Pkg.Id,
		PkgName: pr.Pkg.Name,
		RepoURL: url,
	}
}
