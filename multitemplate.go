package multitemplate

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
)

// Render type
type (
	Render          map[string]*template.Template
	TemplateOptions struct {
		LeftDelimiter  string
		RightDelimiter string
	}
)

type TemplateOption func(*TemplateOptions)

func WithLeftDelimiter(delim string) TemplateOption {
	return func(t *TemplateOptions) {
		t.LeftDelimiter = delim
	}
}

func WithRightDelimiter(delim string) TemplateOption {
	return func(t *TemplateOptions) {
		t.RightDelimiter = delim
	}
}

func Delims(leftDelim, rightDelim string) TemplateOption {
	return func(t *TemplateOptions) {
		WithLeftDelimiter(leftDelim)(t)
		WithRightDelimiter(rightDelim)(t)
	}
}

func NewTemplateOptions(opts ...TemplateOption) *TemplateOptions {
	const (
		defaultLeftDelim  = "{{"
		defaultRightDelim = "}}"
	)

	t := &TemplateOptions{
		LeftDelimiter:  defaultLeftDelim,
		RightDelimiter: defaultRightDelim,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

var (
	_ render.HTMLRender = Render{}
	_ Renderer          = Render{}
)

// New instance
func New() Render {
	return make(Render)
}

// Add new template
func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}
	r[name] = tmpl
}

// AddFromFiles supply add template from files
func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromGlob supply add template from global path
func (r Render) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.ParseGlob(glob))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromFS supply add template from fs.FS (e.g. embed.FS)
func (r Render) AddFromFS(name string, fsys fs.FS, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFS(fsys, files...))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromFSFuncs supply add template from fs.FS (e.g. embed.FS) with callback func
func (r Render) AddFromFSFuncs(name string, funcMap template.FuncMap, fsys fs.FS, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFS(fsys, files...))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromString supply add template from strings
func (r Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New(name).Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromStringsFuncs supply add template from strings
func (r Render) AddFromStringsFuncs(
	name string,
	funcMap template.FuncMap,
	templateStrings ...string,
) *template.Template {
	tmpl := template.New(name).Funcs(funcMap)

	for _, ts := range templateStrings {
		tmpl = template.Must(tmpl.Parse(ts))
	}

	r.Add(name, tmpl)
	return tmpl
}

// AddFromStringsFuncsWithOptions supply add template from strings with options
func (r Render) AddFromStringsFuncsWithOptions(
	name string,
	funcMap template.FuncMap,
	options TemplateOptions,
	templateStrings ...string,
) *template.Template {
	tmpl := template.New(name).
		Delims(options.LeftDelimiter, options.RightDelimiter).
		Funcs(funcMap)

	for _, ts := range templateStrings {
		tmpl = template.Must(
			tmpl.Parse(ts),
		).Delims(options.LeftDelimiter, options.RightDelimiter)
	}

	r.Add(name, tmpl)
	return tmpl
}

// AddFromFilesFuncs supply add template from file callback func
func (r Render) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromFilesFuncsWithOptions supply add template from file callback func with options
func (r Render) AddFromFilesFuncsWithOptions(
	name string,
	funcMap template.FuncMap,
	options TemplateOptions,
	files ...string,
) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(
		template.New(tname).
			Delims(options.LeftDelimiter, options.RightDelimiter).
			Funcs(funcMap).
			ParseFiles(files...),
	)
	r.Add(name, tmpl)
	return tmpl
}

// Instance supply render string
func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r[name],
		Data:     data,
	}
}
