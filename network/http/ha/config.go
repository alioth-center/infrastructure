package ha

type ParamsConfig struct {
	Necessary  []string `json:"necessary,omitempty" yaml:"necessary,omitempty" xml:"necessary,omitempty"`
	Additional []string `json:"additional,omitempty" yaml:"additional,omitempty" xml:"additional,omitempty"`
}

type EndPointConfig struct {
	Name        string       `json:"name" yaml:"name"`
	Path        string       `json:"path" yaml:"path"`
	Methods     []string     `json:"methods" yaml:"methods"`
	Chains      []string     `json:"chains" yaml:"chains"`
	Middlewares []string     `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Headers     ParamsConfig `json:"headers,omitempty" yaml:"headers,omitempty"`
	Paths       ParamsConfig `json:"paths,omitempty" yaml:"paths,omitempty"`
	Queries     ParamsConfig `json:"queries,omitempty" yaml:"queries,omitempty"`
	Cookies     ParamsConfig `json:"cookies,omitempty" yaml:"cookies,omitempty"`
}

type RouterConfig struct {
	Path       string           `json:"path" yaml:"path" xml:"path"`
	SubRouters []RouterConfig   `json:"sub_routers,omitempty" yaml:"sub_routers,omitempty" xml:"sub_routers,omitempty"`
	EndPoints  []EndPointConfig `json:"endpoints,omitempty" yaml:"endpoints,omitempty" xml:"endpoints,omitempty"`
}

type EngineConfig struct {
	Bind    string       `json:"bind" yaml:"bind" xml:"bind"`
	Router  RouterConfig `json:"router" yaml:"router" xml:"router"`
	Serving bool         `json:"serving" yaml:"serving" xml:"serving"`
	Path    string       `json:"path,omitempty" yaml:"path,omitempty" xml:"path,omitempty"`
	Block   bool         `json:"block,omitempty" yaml:"block,omitempty" xml:"block,omitempty"`
}
