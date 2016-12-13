# Multitemplate

[![Build Status](https://travis-ci.org/gin-contrib/multitemplate.svg)](https://travis-ci.org/gin-contrib/multitemplate)
[![codecov](https://codecov.io/gh/gin-contrib/multitemplate/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/multitemplate)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/multitemplate)]`(https://goreportcard.com/report/github.com/gin-contrib/multitemplate)
[![GoDoc](https://godoc.org/github.com/gin-contrib/multitemplate?status.svg)](https://godoc.org/github.com/gin-contrib/multitemplate)

This is a custom HTML render to support multi templates, ie. more than one `*template.Template`.

# Simple example

```go
package main

import (
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/multitemplate"
)

func main() {
    router := gin.Default()
    router.HTMLRender = createMyRender()
    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index", data)
    })
    router.Run(":8080")
}

func createMyRender() multitemplate.Render {
    r := multitemplate.New()
    r.AddFromFiles("index", "base.html", "base.html")
    r.AddFromFiles("article", "base.html", "article.html")
    r.AddFromFiles("login", "base.html", "login.html")
    r.AddFromFiles("dashboard", "base.html", "dashboard.html")

    return r
}
```

## Advanced example

[https://elithrar.github.io/article/approximating-html-template-inheritance/](https://elithrar.github.io/article/approximating-html-template-inheritance/)

```go
package main

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.HTMLRender = loadTemplates("./templates")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "Welcome!",
		})
	})
	router.Run(":8080")
}

func loadTemplates(templatesDir string) multitemplate.Render {
	r := multitemplate.New()

	layouts, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		r.Add(filepath.Base(layout), template.Must(template.ParseFiles(files...)))
	}
	return r
}
```
