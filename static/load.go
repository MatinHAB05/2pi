package static

import (
	"html/template"
)

type StaticFiles struct {
	wakeupHTMLTemplate *template.Template
}

func InitStaticFiles() (*StaticFiles, error) {
	wakeupHTMLTemplate, err := template.ParseFiles("static/wake-up.html")
	if err != nil {
		return nil, err
	}

	return &StaticFiles{
		wakeupHTMLTemplate: wakeupHTMLTemplate,
	}, nil

}
