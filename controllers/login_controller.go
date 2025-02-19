package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"neilz.space/web/services"
)

func RegisterPageRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "register.html",
		gin.H{
			"title":      "Blog Register",
			"isLoggedIn": isLoggedIn,
		})
}
func LoginPageRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")

	c.HTML(http.StatusOK, "login.html",
		gin.H{
			"title":      "Blog Login",
			"isLoggedIn": isLoggedIn,
		})
}
func LoginRoute(c *gin.Context) {
	userID := c.PostForm("id")
	userPassword := c.PostForm("password")
	log.Println(userID, userPassword)
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()
	hostInfo := fmt.Sprintf("%s@%s", userAgent, ipAddress)

	accessToken, refreshToken, err := services.LoginService(userID, userPassword, hostInfo)
	if err != nil {
		ErrorRediect(c, err.Error())
		return
	}
	c.Set("isLoggedIn", true)
	c.SetCookie("access-token", accessToken, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("refresh-token", refreshToken, 60*60*24, "/", "localhost", false, true)
	// c.Header("access-token", accessToken)
	// c.Header("refresh-token", refreshToken)
	c.Redirect(http.StatusSeeOther, "/")
}

func RegisterRoute(c *gin.Context) {
	userID := c.PostForm("id")
	userPassword := c.PostForm("password")
	log.Println(userID, userPassword)
	err := services.RegisterService(userID, userPassword)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	c.Redirect(http.StatusSeeOther, "/")
}
func LogoutRoute(c *gin.Context) {
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()
	hostInfo := fmt.Sprintf("%s@%s", userAgent, ipAddress)

	err := services.LogoutService(hostInfo)
	if err != nil {
		return
	}
	c.SetCookie("access-token", "", 60*60*24, "/", "localhost", false, true)
	c.SetCookie("refresh-token", "", 60*60*24, "/", "localhost", false, true)
	c.Set("isLoggedIn", false)
	c.Redirect(http.StatusSeeOther, "/")
}
