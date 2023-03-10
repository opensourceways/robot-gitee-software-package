package emailimpl

type Config struct {
	EmailServer     EmailServer `json:"email"`
	MaintainerEmail string      `json:"maintainer_email" require:"true"`
}

type EmailServer struct {
	AuthCode string `json:"auth_code" required:"true"`
	From     string `json:"from"      required:"true"`
	Host     string `json:"host"      required:"true"`
	Port     int    `json:"port"      required:"true"`
}

func (cfg *Config) SetDefault() {

}

func (cfg *Config) Validate() error {
	return nil
}
