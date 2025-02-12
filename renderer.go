package multitemplate

import (
	"html/template"
	"io/fs"

	"github.com/gin-gonic/gin/render"
)

// Renderer type is the Agnostic Renderer for multitemplates.
// When gin is in debug mode then all multitemplates works with
// hot reloading allowing you modify file templates and seeing changes instantly.
// Renderer should be created using multitemplate.NewRenderer() constructor.
type Renderer interface {
	render.HTMLRender
	Add(name string, tmpl *template.Template)
	AddFromFiles(name string, files ...string) *template.Template
	AddFromGlob(name, glob string) *template.Template
	AddFromFS(name string, fsys fs.FS, files ...string) *template.Template
	AddFromFSFuncs(name string, funcMap template.FuncMap, fsys fs.FS, files ...string) *template.Template
	AddFromString(name, templateString string) *template.Template
	AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template
	AddFromStringsFuncsWithOptions(
		name string,
		funcMap template.FuncMap,
		options TemplateOptions,
		templateStrings ...string,
	) *template.Template
	AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template
	AddFromFilesFuncsWithOptions(
		name string,
		funcMap template.FuncMap,
		options TemplateOptions,
		files ...string,
	) *template.Template
}
