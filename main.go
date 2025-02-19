package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"neilz.space/web/controllers"
	"neilz.space/web/middlewares"
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

	r.Use(middlewares.CheckLoginned)

	r.GET("/", controllers.IndexRoute)
	r.GET("/about", controllers.AboutRoute)

	r.GET("/error", controllers.ErrorRoute)

	r.GET("/blog/:pageNumber", controllers.BlogListRoute)
	r.GET("/blog-article/:articleNumber", controllers.BlogArticleRoute)
	// r.GET("/blog-post", controllers.BlogPostPageRoute)
	// r.POST("/blog-posting", controllers.BlogPostingRoute)
	r.GET("/blog-post-test", controllers.AllBlogListJSON)

	r.GET("/login", controllers.LoginPageRoute)
	r.GET("/logout", controllers.LogoutRoute)

	r.POST("/logining", controllers.LoginRoute)

	r.GET("/register", controllers.RegisterPageRoute)
	r.POST("/registering", controllers.RegisterRoute)

	r.GET("/opensource", controllers.OpensourcePageRoute)

	auth := r.Group("/auth")
	auth.Use(middlewares.RequireAuth)
	auth.POST("/blog-posting", controllers.BlogPostingRoute)
	auth.GET("/blog-post", controllers.BlogPostPageRoute)
	auth.GET("/blog-edit/:articleNumber", controllers.BlogEditPageRoute)
	auth.POST("/blog-editing/:articleNumber", controllers.BlogEditingRoute)
	auth.GET("/blog-remove/:articleNumber", controllers.BlogRemoveRoute)
	auth.POST("/blog-removing/:articleNumber", controllers.BlogRemovingRoute)
	r.Run(":8080")
}
