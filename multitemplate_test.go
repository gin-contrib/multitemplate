package multitemplate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createFromFile() Render {
	r := New()
	r.AddFromFiles("index", "tests/base.html", "tests/article.html")

	return r
}

func createFromGlob() Render {
	r := New()
	r.AddFromGlob("index", "tests/global/*")

	return r
}

func createFromString() Render {
	r := New()
	r.AddFromString("index", "Welcome to {{ .name }} template")

	return r
}

func TestAddFromFiles(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromFile()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Test Multiple Template",
		})
	})

	w := performRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>Test Multiple Template</p>\nHi, this is article template\n", w.Body.String())
}

func TestAddFromGlob(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromGlob()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Test Multiple Template",
		})
	})

	w := performRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>Test Multiple Template</p>\nHi, this is login template\n", w.Body.String())
}

func TestAddFromString(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromString()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"name": "index",
		})
	})

	w := performRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome to index template", w.Body.String())
}
