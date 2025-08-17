package types

type Config struct {
	Dev   configDev   `json:"dev"`
	Build configBuild `json:"build"`
}

type configDev struct {
	Port uint16 `json:"port"`
}

type configBuild struct {
	Directory string          `json:"directory"`
	Html      configBuildHtml `json:"html"`
	Css       configBuildCss  `json:"css"`
	Js        configBuildJs   `json:"js"`
}

type configBuildHtml struct {
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

type configBuildCss struct {
	KeepCSS2  bool `json:"keepCSS2"`
	Precision int  `json:"precision"`
	Inline    bool `json:"inline"`
}

type configBuildJs struct {
	Precision    int  `json:"precision"`
	KeepVarNames bool `json:"keepVarNames"`
	Version      int  `json:"version"`
}
