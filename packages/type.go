package packages

type Settings struct {
	ProgramDir   string   `json:"programDir"`
	SubscribeURL string   `json:"subscribeURL"`
	Base64File   string   `json:"base64File"`
	TempListPath string   `json:"tempListPath"`
	TempJsonPath string   `json:"tempJsonPath"`
	TemplatePath string   `json:"templatePath"`
	Filter       []Filter `json:"filter"`
	OutputPath   string   `json:"outputPath"`
	StartStep    int      `json:"startStep"`
}

type Filter struct {
	Action   string   `json:"action"`
	Keywords []string `json:"keywords"`
}

type TrojanConfig struct {
	Tag        string    `json:"tag"`
	Type       string    `json:"type"`
	Server     string    `json:"server"`
	ServerPort int       `json:"server_port"`
	Password   string    `json:"password"`
	TLS        trojanTLS `json:"tls"`
}
type trojanTLS struct {
	Enable     bool   `json:"enabled"`
	Insecure   bool   `json:"insecure"`
	ServerName string `json:"server_name"`
}

type ShadowsocksConfig struct {
	Tag        string `json:"tag"`
	Type       string `json:"type"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
}
