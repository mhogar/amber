package renderer

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
)

type TemplateData struct {
	AppName string
	BaseURL string
	Data    interface{}
}

type Renderer interface {
	// RenderView renders the provided templates with provided data.
	// Returns the raw rendered view and panics on error.
	RenderView(req *http.Request, data interface{}, templates ...string) []byte
}

type CoreRenderer struct {
	Templates map[string]*template.Template
}

func (r CoreRenderer) RenderView(req *http.Request, data interface{}, templates ...string) []byte {
	//create the data object
	d := TemplateData{
		AppName: config.GetAppName(),
		BaseURL: "http://" + req.Host,
		Data:    data,
	}

	//update the template paths
	templates = append(templates, "partials/base")
	for index, t := range templates {
		templates[index] = config.GetAppRoot("views", t+".gohtml")
	}

	//fetch or parse the template
	t, ok := r.Templates[templates[0]]
	if !ok {
		t = template.Must(template.ParseFiles(templates...))
		r.Templates[templates[0]] = t
	}

	//render the template
	var buffer bytes.Buffer
	err := t.Execute(&buffer, d)
	if err != nil {
		log.Panicf(common.ChainError("error rendering template(s)", err).Error())
	}

	return buffer.Bytes()
}
