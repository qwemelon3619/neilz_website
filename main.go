package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"neilz.space/web/controllers"
	"neilz.space/web/models"
	"neilz.space/web/setting"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	// handle err
	time.Local = loc // -> this is setting the global timezone

	models.ConnectDataBase()
	defer models.CloseDataBase()

	r := gin.Default()
	r.Static("/css", "./templates/css")
	r.Static("/js", "./templates/js")
	r.Static("/assets", "./templates/assets")

	setting.AddTemplateFunction(r)
	r.LoadHTMLGlob("templates/html/*")

	r.GET("/", controllers.IndexRoute)
	r.GET("/about", controllers.AboutRoute)
	// r.GET("/contact", controllers.ContactRoute)
	// r.GET("/blog", controllers.BlogHomeRoute)

	r.GET("/blog/:pageNumber", controllers.BlogListRoute)
	r.GET("/blog-article/:articleNumber", controllers.BlogArticleRoute)
	r.GET("/error", controllers.ErrorRoute)
	r.GET("/blog-post", controllers.BlogPostPageRoute)
	r.POST("/blog-posting", controllers.BlogPostingRoute)

	r.GET("/blog-post-test", controllers.AllBlogListJSON)
	// r.NoRoute(func(c *gin.Context) { controllers.ErrorRediect(c, "wrong url") })
	r.Run(":8080")
}
