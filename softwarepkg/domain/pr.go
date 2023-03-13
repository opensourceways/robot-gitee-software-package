package domain

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

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}

func ToSoftwarePkgRepo(pr *PullRequest) *SoftwarePkgRepo {
	return &SoftwarePkgRepo{
		Pkg: pr.Pkg,
	}
}

func (p *PullRequest) IsMerged() bool {
	return p.Merged
}
