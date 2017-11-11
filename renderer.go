package multitemplate

import "html/template"
import "github.com/gin-gonic/gin/render"

// Renderer type
type Renderer interface {
	render.HTMLRender
	Add(name string, tmpl *template.Template)
	AddFromFiles(name string, files ...string) *template.Template
	AddFromGlob(name, glob string) *template.Template
	AddFromString(name, templateString string) *template.Template
	AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template
	AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template
}
