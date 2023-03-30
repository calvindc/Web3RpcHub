package handlers

import "html/template"

type changeLanguageTemplateData struct {
	PostRoute    string
	CSRFElement  template.HTML
	LangTag      string
	RedirectPage string
	Translation  string
	ClassList    string
}

var changeLanguageTemplate = template.Must(template.New("changeLanguageForm").Parse(`
foo...`))
