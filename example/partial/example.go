package main

import (
	"log"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/item.html")
	return r
}

func main() {
	router := gin.Default()
	router.HTMLRender = createMyRender()

	// Route to render full template
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"items": []gin.H{
				{"name": "Apple"},
				{"name": "Banana"},
				{"name": "Cherry"},
			},
		})
	})

	// Route to render partial template using "index#item" syntax
	router.GET("/item", func(c *gin.Context) {
		c.HTML(200, "index#item", gin.H{
			"name": "Watermelon",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
