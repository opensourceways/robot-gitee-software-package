package repositoryimpl

import (
	"fmt"
	"os"

	"github.com/opensourceways/server-common-lib/utils"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

const prDir = "pr-storage"

func NewRepositoryImpl() *RepositoryImpl {
	return &RepositoryImpl{}
}

type RepositoryImpl struct {
}

func (impl *RepositoryImpl) Add(pr *domain.PullRequest) error {
	data, err := yaml.Marshal(pr)
	if err != nil {
		return err
	}

	fileName, err := impl.genFileName(pr.Num)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func (impl *RepositoryImpl) Find(prNum int) (pr domain.PullRequest, err error) {
	fileName, err := impl.genFileName(prNum)
	if err != nil {
		return
	}

	_, err = os.Stat(fileName)
	if err != nil {
		return
	}

	if err = utils.LoadFromYaml(fileName, &pr); err != nil {
		return
	}

	return
}

func (impl *RepositoryImpl) genFileName(prNum int) (string, error) {
	if s, err := os.Stat(prDir); err != nil || !s.IsDir() {
		if err = os.Mkdir(prDir, 755); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/%d.yaml", prDir, prNum), nil
}
