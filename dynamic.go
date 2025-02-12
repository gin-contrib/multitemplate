package multitemplate

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// DynamicRender type
type DynamicRender map[string]*templateBuilder

var (
	_ render.HTMLRender = DynamicRender{}
	_ Renderer          = DynamicRender{}
)

// NewDynamic is the constructor for Dynamic templates
func NewDynamic() DynamicRender {
	return make(DynamicRender)
}

// NewRenderer allows create an agnostic multitemplate renderer
// depending on enabled gin mode
func NewRenderer() Renderer {
	if gin.IsDebugging() {
		return NewDynamic()
	}
	return New()
}

// Type of dynamic builder
type builderType int

// Types of dynamic builders
const (
	templateType builderType = iota
	filesTemplateType
	globTemplateType
	fsTemplateType
	fsFuncTemplateType
	stringTemplateType
	stringFuncTemplateType
	filesFuncTemplateType
)

// Builder for dynamic templates
type templateBuilder struct {
	buildType       builderType
	tmpl            *template.Template
	templateName    string
	files           []string
	glob            string
	fsys            fs.FS
	templateString  string
	funcMap         template.FuncMap
	templateStrings []string
	options         TemplateOptions
}

func (tb templateBuilder) buildTemplate() *template.Template {
	switch tb.buildType {
	case templateType:
		return tb.tmpl.Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter)
	case filesTemplateType:
		tmpl := template.Must(template.ParseFiles(tb.files...))
		return tmpl.Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter)
	case globTemplateType:
		tmpl := template.Must(template.ParseGlob(tb.glob))
		return tmpl.Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter)
	case fsTemplateType:
		tmpl := template.Must(template.ParseFS(tb.fsys, tb.files...))
		return tmpl.Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter)
	case fsFuncTemplateType:
		tmpl := template.New(tb.templateName).
			Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter).
			Funcs(tb.funcMap)
		return template.Must(tmpl.ParseFS(tb.fsys, tb.files...))
	case stringTemplateType:
		tmpl := template.New(tb.templateName).
			Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter)
		return template.Must(tmpl.Parse(tb.templateString))
	case stringFuncTemplateType:
		tmpl := template.New(tb.templateName).
			Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter).
			Funcs(tb.funcMap)
		for _, ts := range tb.templateStrings {
			tmpl = template.Must(tmpl.Parse(ts))
		}
		return tmpl
	case filesFuncTemplateType:
		tmpl := template.New(tb.templateName).
			Delims(tb.options.LeftDelimiter, tb.options.RightDelimiter).
			Funcs(tb.funcMap)
		return template.Must(tmpl.ParseFiles(tb.files...))
	default:
		panic("Invalid builder type for dynamic template")
	}
}

// Add new template
func (r DynamicRender) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template cannot be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	builder := &templateBuilder{templateName: name, tmpl: tmpl, options: *NewTemplateOptions()}
	builder.buildType = templateType
	r[name] = builder
}

// AddFromFiles supply add template from files
func (r DynamicRender) AddFromFiles(name string, files ...string) *template.Template {
	builder := &templateBuilder{templateName: name, files: files, options: *NewTemplateOptions()}
	builder.buildType = filesTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromGlob supply add template from global path
func (r DynamicRender) AddFromGlob(name, glob string) *template.Template {
	builder := &templateBuilder{templateName: name, glob: glob, options: *NewTemplateOptions()}
	builder.buildType = globTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromFS adds a new template to the DynamicRender from the provided file system (fs.FS) and files.
// It allows you to specify a custom function map (funcMap) to be used within the template.
// The name parameter is used to associate the template with a key in the DynamicRender.
// The files parameter is a variadic list of file paths to be included in the template.
//   - name: The name to associate with the template in the DynamicRender.
//   - fsys: The file system (fs.FS) from which to read the template files.
//   - files: A variadic list of file paths to be included in the template.
//
// Returns:
//   - *template.Template: The constructed template.
func (r DynamicRender) AddFromFS(name string, fsys fs.FS, files ...string) *template.Template {
	builder := &templateBuilder{templateName: name, fsys: fsys, files: files}
	builder.buildType = fsTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromFSFuncs adds a new template to the DynamicRender from the provided file system (fs.FS) and files.
// It allows you to specify a custom function map (funcMap) to be used within the template.
//
// Parameters:
//   - name: The name to associate with the template in the DynamicRender.
//   - funcMap: A map of functions to be used within the template.
//   - fsys: The file system (fs.FS) from which to read the template files.
//   - files: A variadic list of file paths to be included in the template.
//
// Returns:
//   - *template.Template: The constructed template.
func (r DynamicRender) AddFromFSFuncs(
	name string,
	funcMap template.FuncMap,
	fsys fs.FS,
	files ...string,
) *template.Template {
	tname := filepath.Base(files[0])
	builder := &templateBuilder{
		templateName: tname,
		funcMap:      funcMap,
		fsys:         fsys,
		files:        files,
	}
	builder.buildType = fsFuncTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromString supply add template from strings
func (r DynamicRender) AddFromString(name, templateString string) *template.Template {
	builder := &templateBuilder{templateName: name, templateString: templateString, options: *NewTemplateOptions()}
	builder.buildType = stringTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromStringsFuncs supply add template from strings
func (r DynamicRender) AddFromStringsFuncs(
	name string,
	funcMap template.FuncMap,
	templateStrings ...string,
) *template.Template {
	builder := &templateBuilder{
		templateName: name, funcMap: funcMap,
		templateStrings: templateStrings,
		options:         *NewTemplateOptions(),
	}
	builder.buildType = stringFuncTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromStringsFuncsWithOptions supply add template from strings with options
func (r DynamicRender) AddFromStringsFuncsWithOptions(
	name string,
	funcMap template.FuncMap,
	options TemplateOptions,
	templateStrings ...string,
) *template.Template {
	builder := &templateBuilder{
		templateName:    name,
		funcMap:         funcMap,
		templateStrings: templateStrings,
		options:         options,
	}
	builder.buildType = stringFuncTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromFilesFuncs supply add template from file callback func
func (r DynamicRender) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	builder := &templateBuilder{templateName: tname, funcMap: funcMap, files: files, options: *NewTemplateOptions()}
	builder.buildType = filesFuncTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromFilesFuncs supply add template from file callback func
func (r DynamicRender) AddFromFilesFuncsWithOptions(
	name string,
	funcMap template.FuncMap,
	options TemplateOptions,
	files ...string,
) *template.Template {
	tname := filepath.Base(files[0])
	builder := &templateBuilder{
		templateName: tname,
		funcMap:      funcMap,
		files:        files,
		options:      options,
	}
	builder.buildType = filesFuncTemplateType
	r[name] = builder
	return builder.buildTemplate()
}

// Instance supply render string
func (r DynamicRender) Instance(name string, data interface{}) render.Render {
	builder, ok := r[name]
	if !ok {
		panic(fmt.Sprintf("Dynamic template with name %s not found", name))
	}
	return render.HTML{
		Template: builder.buildTemplate(),
		Data:     data,
	}
}
