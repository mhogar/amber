package renderer

import (
	"authserver/config"
	"bytes"
	"log"
	"path"
	"text/template"
)

type Renderer interface {
	// RenderView renders the view with the given name and data.
	// Returns the raw rendered view, and panics on error.
	RenderView(name string, data interface{}) []byte
}

type CoreRenderer struct{}

func (CoreRenderer) RenderView(name string, data interface{}) []byte {
	//parse the template
	t := template.Must(template.ParseFiles(path.Join(config.GetAppRoot(), "views", name)))

	//render the template
	var buffer bytes.Buffer
	err := t.Execute(&buffer, data)
	if err != nil {
		log.Panicf("error rendering %s template", name)
	}

	return buffer.Bytes()
}
