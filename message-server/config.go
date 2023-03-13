package messageserver

type Config struct {
	GroupName string `json:"group_name"    required:"true"`
	Topics    Topics `json:"topics"        required:"true"`
}

type Topics struct {
	CIPassed            string `json:"ci_passed"                    required:"true"`
	ApplyingSoftwarePkg string `json:"applying_software_pkg"        required:"true"`
	ApprovedSoftwarePkg string `json:"approved_software_pkg"        required:"true"`
}
