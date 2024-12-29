package packages

type Settings struct {
	ProgramDir   string   `json:"programDir"`
	SubscribeURL string   `json:"subscribeURL"`
	SubURLs      []SubURL `json:"subURLs"`
	Base64File   string   `json:"base64File"`
	TempListPath string   `json:"tempListPath"`
	TempJsonPath string   `json:"tempJsonPath"`
	TemplatePath string   `json:"templatePath"`
	Filter       []Filter `json:"filter"`
	OutputPath   string   `json:"outputPath"`
	StartStep    int      `json:"startStep"`
}

type SubURL struct {
	Tag string   `json:"tag"`
	URL string `json:"url"`
}

type Filter struct {
	Action   string   `json:"action"`
	Keywords []string `json:"keywords"`
}
