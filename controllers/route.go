package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func IndexRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "HomePage"})
}

func AboutRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{"title": "About"})
}

func ContactRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", gin.H{"title": "Contact"})
}
func Redirect(c *gin.Context, path string) {
	c.HTML(http.StatusOK, "redirect.html", gin.H{
		"RedirectURL": path,
	})
}
func ErrorRoute(c *gin.Context) {
	errorMessage := c.Query("error")
	c.HTML(http.StatusOK, "error.html", gin.H{"title": "Error", "error": errorMessage})
}
func ErrorRediect(c *gin.Context, errorMessage string) {
	log.Println(errorMessage)
	c.Redirect(http.StatusSeeOther, "/error?error="+url.QueryEscape(errorMessage))
}
