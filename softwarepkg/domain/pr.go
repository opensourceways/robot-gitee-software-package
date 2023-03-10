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
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SoftwarePkg struct {
	SoftwarePkgBasic

	ImporterName  string
	ImporterEmail string
	Application   SoftwarePkgApplication
}

// PullRequest
type PullRequest struct {
	Num  int              `json:"num"`
	Link string           `json:"link"`
	Pkg  SoftwarePkgBasic `json:"pkg"`
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}
