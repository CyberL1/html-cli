package utils

import (
	"html-cli/constants"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func Minify(fileBytes []byte) []byte {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags:        constants.Config.Build.Html.KeepDocumentTags,
		KeepComments:            constants.Config.Build.Html.KeepComments,
		KeepConditionalComments: constants.Config.Build.Html.KeepConditionalComments,
		KeepSpecialComments:     constants.Config.Build.Html.KeepSpecialComments,
		KeepDefaultAttrVals:     constants.Config.Build.Html.KeepDefaultAttrVals,
		KeepEndTags:             constants.Config.Build.Html.KeepEndTags,
		KeepQuotes:              constants.Config.Build.Html.KeepQuotes,
		KeepWhitespace:          constants.Config.Build.Html.KeepWhitespace,
		TemplateDelims:          constants.Config.Build.Html.TemplateDelims,
	})

	m.Add("text/css", &css.Minifier{
		KeepCSS2:  constants.Config.Build.Css.KeepCSS2,
		Precision: constants.Config.Build.Css.Precision,
		Inline:    constants.Config.Build.Css.Inline,
	})

	m.AddRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), &js.Minifier{
		Precision:    constants.Config.Build.Js.Precision,
		KeepVarNames: constants.Config.Build.Js.KeepVarNames,
		Version:      constants.Config.Build.Js.Version,
	})

	minified, err := m.String("text/html", string(fileBytes))
	if err != nil {
		panic(err)
	}
	return []byte(minified)
}
