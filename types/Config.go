package types

type Config struct {
	Dev   ConfigDev   `json:"dev"`
	Build ConfigBuild `json:"build"`
}

type ConfigDev struct {
	Port uint16 `json:"port"`
}

type ConfigBuild struct {
	Directory string          `json:"directory"`
	Html      ConfigBuildHtml `json:"html"`
	Css       ConfigBuildCss  `json:"css"`
	Js        ConfigBuildJs   `json:"js"`
}

type ConfigBuildHtml struct {
	KeepDocumentTags        bool      `json:"keepDocumentTags"`
	KeepComments            bool      `json:"keepComments"`
	KeepConditionalComments bool      `json:"keepConditionalComments"`
	KeepSpecialComments     bool      `json:"keepSpecialComments"`
	KeepDefaultAttrVals     bool      `json:"keepDefaultAttrVals"`
	KeepEndTags             bool      `json:"keepEndTags"`
	KeepQuotes              bool      `json:"keepQuotes"`
	KeepWhitespace          bool      `json:"keepWhitespace"`
	TemplateDelims          [2]string `json:"templateDelims"`
}

type ConfigBuildCss struct {
	KeepCSS2  bool `json:"keepCSS2"`
	Precision int  `json:"precision"`
	Inline    bool `json:"inline"`
}

type ConfigBuildJs struct {
	Precision    int  `json:"precision"`
	KeepVarNames bool `json:"keepVarNames"`
	Version      int  `json:"version"`
}
