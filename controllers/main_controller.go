package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"neilz.space/web/models"
)

func IndexRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	articles, err := models.GetAllBlogArticles()
	if err != nil {
		ErrorRediect(c, "DB error")
	}
	if len(articles) >= 3 {
		articles = articles[:3]
	}
	c.HTML(http.StatusOK, "index.html",
		gin.H{
			"title":      "HomePage",
			"isLoggedIn": isLoggedIn,
			"articles":   articles,
		})
}

func AboutRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "about.html",
		gin.H{
			"title":      "About",
			"isLoggedIn": isLoggedIn,
		})
}

func ContactRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "contact.html",
		gin.H{
			"title":      "Contact",
			"isLoggedIn": isLoggedIn,
		})

}
func ErrorRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	errorMessage := c.Query("error")
	c.HTML(http.StatusOK, "error.html", gin.H{"title": "Error", "error": errorMessage, "isLoggedIn": isLoggedIn})
}
func ErrorRediect(c *gin.Context, errorMessage string) {
	log.Println(errorMessage)
	c.Redirect(http.StatusSeeOther, "/error?error="+url.QueryEscape(errorMessage))
}

func OpensourcePageRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "opensource.html",
		gin.H{
			"title": "Open Source Lisence",
		})
}
