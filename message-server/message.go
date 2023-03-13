package messageserver

import (
	"regexp"
	"strconv"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

var prNumRe = regexp.MustCompile(`\d+$`)

type messageOfNewPkg struct{}

func (msg *messageOfNewPkg) toCmd() (app.CmdToCreatePR, error) {
	return app.CmdToCreatePR{}, nil
}

type messageOfApprovedPkg struct {
	PkgId      string `json:"pkg_id"`
	PkgName    string `json:"pkg_name"`
	RelevantPR string `json:"pr"`
}

func (msg *messageOfApprovedPkg) toCmd() (cmd app.CmdToMergePR, err error) {
	prNumStr := prNumRe.FindString(msg.RelevantPR)
	prNumInt, err := strconv.Atoi(prNumStr)
	if err != nil {
		return
	}

	cmd.PRNum = prNumInt

	return
}
