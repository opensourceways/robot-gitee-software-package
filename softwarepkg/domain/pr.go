package domain

import "encoding/json"

type SoftwarePkgSourceCode struct {
	Address string
	License string
}

type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageDesc       string
	PackagePlatform   string
	ImportingPkgSig   string
	ReasonToImportPkg string
}

type SoftwarePkgBasic struct {
	Id   string
	Name string
}

type SoftwarePkg struct {
	SoftwarePkgBasic

	ImporterName  string
	ImporterEmail string
	Application   SoftwarePkgApplication
}

// PullRequest
type PullRequest struct {
	Num    int
	Link   string
	Merged bool
	Pkg    SoftwarePkgBasic
}

func (r *PullRequest) SetMerged() {
	r.Merged = true
}

func (r *PullRequest) IsMerged() bool {
	return r.Merged
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}

func ToSoftwarePkgRepo(pr *PullRequest, url string) *SoftwarePkgRepo {
	return &SoftwarePkgRepo{
		Pkg:     pr.Pkg,
		RepoURL: url,
	}
}

func (s *SoftwarePkgRepo) Message() ([]byte, error) {
	return json.Marshal(s)
}
