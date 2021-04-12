package conf

type Config struct {
	Discover DiscoverConf `yaml:"discover"`
	Rpcx     RpcxConf     `yaml:"rpcx"`
	Stores   StoresConf   `yaml:"stores"`
	Coder    CoderConf    `yaml:"coder"`
}

type DiscoverConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
}

type RpcxConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
}

type StoresConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
	Password    string   `yaml:"password"`
	DB          int      `yaml:"db"`
}

type CoderConf struct {
	UseType string `yaml:"use_type"`
}
