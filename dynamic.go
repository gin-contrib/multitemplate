package multitemplate

import (
	"html/template"

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

// NewMultitemplateRenderer allows create an agnostic multitemplate renderer
// depending on enabled gin mode
func NewMultitemplateRenderer() Renderer {
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
	templateString  string
	funcMap         template.FuncMap
	templateStrings []string
}

func (tb templateBuilder) buildTemplate() *template.Template {
	switch tb.buildType {
	case templateType:
		return tb.tmpl
	case filesTemplateType:
		return template.Must(template.ParseFiles(tb.files...))
	case globTemplateType:
		return template.Must(template.ParseGlob(tb.glob))
	case stringTemplateType:
		return template.Must(template.New(tb.templateName).Parse(tb.templateString))
	case stringFuncTemplateType:
		tmpl := template.New(tb.templateName).Funcs(tb.funcMap)
		for _, ts := range tb.templateStrings {
			tmpl = template.Must(tmpl.Parse(ts))
		}
		return tmpl
	case filesFuncTemplateType:
		return template.Must(template.New(tb.templateName).Funcs(tb.funcMap).ParseFiles(tb.files...))
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
	builder := &templateBuilder{templateName: name, tmpl: tmpl}
	r[name] = builder
}

// AddFromFiles supply add template from files
func (r DynamicRender) AddFromFiles(name string, files ...string) *template.Template {
	builder := &templateBuilder{templateName: name, files: files}
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromGlob supply add template from global path
func (r DynamicRender) AddFromGlob(name, glob string) *template.Template {
	builder := &templateBuilder{templateName: name, glob: glob}
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromString supply add template from strings
func (r DynamicRender) AddFromString(name, templateString string) *template.Template {
	builder := &templateBuilder{templateName: name, templateString: templateString}
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromStringsFuncs supply add template from strings
func (r DynamicRender) AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template {
	builder := &templateBuilder{templateName: name, funcMap: funcMap,
		templateStrings: templateStrings}
	r[name] = builder
	return builder.buildTemplate()
}

// AddFromFilesFuncs supply add template from file callback func
func (r DynamicRender) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	builder := &templateBuilder{templateName: name, funcMap: funcMap, files: files}
	r[name] = builder
	return builder.buildTemplate()
}

// Instance supply render string
func (r DynamicRender) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r[name].buildTemplate(),
		Data:     data,
	}
}