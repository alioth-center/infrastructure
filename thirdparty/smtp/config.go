package smtp

type Config struct {
	Secret    string `yaml:"secret,omitempty" json:"secret,omitempty" xml:"secret,omitempty"`
	Server    string `yaml:"server,omitempty" json:"server,omitempty" xml:"server,omitempty"`
	Port      int    `yaml:"port,omitempty" json:"port,omitempty" xml:"port,omitempty"`
	Sender    string `yaml:"sender,omitempty" json:"sender,omitempty" xml:"sender,omitempty"`
	Signature string `yaml:"signature,omitempty" json:"signature,omitempty" xml:"signature,omitempty"`
	EnableTLS bool   `yaml:"enable_tls,omitempty" json:"enable_tls,omitempty" xml:"enable_tls,omitempty"`
}
