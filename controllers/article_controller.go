package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"neilz.space/web/models"
	"neilz.space/web/utils"
)

func BlogArticleRoute(c *gin.Context) {
	articleNumber_s := c.Param("articleNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	article, err := models.GetBlogArticle(articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	content := template.HTML(article.Content)
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-article.html",
		gin.H{
			"title":         article.Title,
			"blogArticle":   article,
			"content":       content,
			"isLoggedIn":    isLoggedIn,
			"articleNumber": articleNumber,
		})
}
func BlogPostPageRoute(c *gin.Context) {
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-post.html",
		gin.H{
			"title":      "Blog Posting",
			"isLoggedIn": isLoggedIn,
		})
}

func BlogPostingRoute(c *gin.Context) {
	article := models.BlogArticleInput{}
	article.Title = c.PostForm("article-title")
	article.SubTitle = c.PostForm("article-subtitle")
	// article.Content = c.PostForm("article-content")

	content := c.PostForm("article-content")
	removeContentString := []string{"<select class=\"ql-ui\" contenteditable=\"false\"><option value=\"plain\">Plain</option><option value=\"bash\">Bash</option><option value=\"cpp\">C++</option><option value=\"cs\">C#</option><option value=\"css\">CSS</option><option value=\"diff\">Diff</option><option value=\"xml\">HTML/XML</option><option value=\"java\">Java</option><option value=\"javascript\">JavaScript</option><option value=\"markdown\">Markdown</option><option value=\"php\">PHP</option><option value=\"python\">Python</option><option value=\"ruby\">Ruby</option><option value=\"sql\">SQL</option></select>", "contenteditable=\"true\""}
	for _, value := range removeContentString {
		content = strings.Replace(content, value, "", -1)
	}
	article.Content = content

	img, err := c.FormFile("article-img")
	if err != nil && err != http.ErrMissingFile {
		ErrorRediect(c, err.Error())
		return
	} else if err == http.ErrMissingFile {
		article.HeaderImageName = "sample.jpg"
	} else if img != nil {
		timestamp := time.Now().UnixNano()
		randomString := utils.GenerateRandomString(8)
		filename_e := strings.Split(img.Filename, ".")
		file_e := filename_e[len(filename_e)-1]
		img.Filename = strconv.Itoa(int(timestamp)) + "_" + randomString + "." + file_e
		if err != nil {
			ErrorRediect(c, err.Error())
			return
		}
		filePath := filepath.Join("./templates/assets/blog_img/", img.Filename)
		log.Print(filePath)
		article.HeaderImageName = img.Filename

		err = c.SaveUploadedFile(img, filePath)
		if err != nil {
			ErrorRediect(c, err.Error())
			return
		}

		err = utils.ImageResize(filePath) //crop image
		if err != nil {
			ErrorRediect(c, err.Error())
			return
		}
	}

	err = models.SaveBlogArticle(article)
	if err != nil {
		ErrorRediect(c, err.Error())
		return
	}

	// Redirect(c, "/blog/1")
	c.Redirect(http.StatusSeeOther, "/blog/1")
}

type PageData struct {
	PageNumber   int
	TotalPages   int
	VisiblePages []int
}

func BlogListRoute(c *gin.Context) {
	pageNumber_S := c.Param("pageNumber")
	var err error
	pageNumber, err := strconv.Atoi(pageNumber_S)
	if err != nil {
		ErrorRediect(c, "Bad request")
	}

	blogArticles, err := models.GetPageBlogArticles(pageNumber)
	if err != nil {
		ErrorRediect(c, "Bad request")
	}

	numberofBlogArticles, err := models.GetNumberOfBlogArticles()
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	totalPages := numberofBlogArticles/6 + 1
	startpage := pageNumber - 2
	endpage := pageNumber + 2
	if startpage <= 1 {
		startpage = 1
	}
	if endpage >= totalPages {
		endpage = totalPages
	}
	var visiblePages []int
	for i := startpage; i <= endpage; i++ {
		visiblePages = append(visiblePages, i)
	}
	pagination := PageData{
		PageNumber:   pageNumber, // 실제 페이지 번호
		TotalPages:   totalPages, // 실제 전체 페이지 수
		VisiblePages: visiblePages,
	}
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-list.html", gin.H{
		"title":        "Blog List",
		"blogArticles": blogArticles,
		"pagination":   pagination,
		"isLoggedIn":   isLoggedIn,
	})
}
func AllBlogListJSON(c *gin.Context) {
	articles, err := models.GetAllBlogArticles()
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	c.JSON(http.StatusOK, articles)
}

func BlogEditPageRoute(c *gin.Context) {
	articleNumber_s := c.Param("articleNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	article, err := models.GetBlogArticle(articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	content := template.HTML(article.Content)
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-edit.html",
		gin.H{
			"title":         "Blog Editing",
			"isLoggedIn":    isLoggedIn,
			"content":       content,
			"article":       article,
			"articleNumber": articleNumber,
		})
}

func BlogEditingRoute(c *gin.Context) {
	article := models.BlogArticleInput{}
	articleNumber_s := c.Param("articleNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
		return
	}
	article.Title = c.PostForm("article-title")
	article.SubTitle = c.PostForm("article-subtitle")
	content := c.PostForm("article-content")
	removeContentString := []string{"<select class=\"ql-ui\" contenteditable=\"false\"><option value=\"plain\">Plain</option><option value=\"bash\">Bash</option><option value=\"cpp\">C++</option><option value=\"cs\">C#</option><option value=\"css\">CSS</option><option value=\"diff\">Diff</option><option value=\"xml\">HTML/XML</option><option value=\"java\">Java</option><option value=\"javascript\">JavaScript</option><option value=\"markdown\">Markdown</option><option value=\"php\">PHP</option><option value=\"python\">Python</option><option value=\"ruby\">Ruby</option><option value=\"sql\">SQL</option></select>", "contenteditable=\"true\""}
	for _, value := range removeContentString {
		content = strings.Replace(content, value, "", -1)
	}
	article.Content = content

	img, err := c.FormFile("article-img")
	if err != nil && err != http.ErrMissingFile {
		ErrorRediect(c, err.Error())
		return
	}
	if img != nil {
		filePath := filepath.Join("./templates/assets/blog_img/", img.Filename)
		log.Print(filePath)
		article.HeaderImageName = img.Filename

		err = c.SaveUploadedFile(img, filePath)
		if err != nil {
			ErrorRediect(c, err.Error())
			return
		}

		err = utils.ImageResize(filePath) //crop image
		if err != nil {
			ErrorRediect(c, err.Error())
			return
		}
	}
	err = models.EditBlogArticle(article, articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
		return
	}
	// Redirect(c, "/blog/1")
	c.Redirect(http.StatusSeeOther, "/blog/1")
}
func BlogRemoveRoute(c *gin.Context) {
	articleNumber_s := c.Param("articleNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	article, err := models.GetBlogArticle(articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	content := template.HTML(article.Content)
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-remove.html",
		gin.H{
			"title":         article.Title,
			"blogArticle":   article,
			"content":       content,
			"isLoggedIn":    isLoggedIn,
			"articleNumber": articleNumber,
		})
}
func BlogRemovingRoute(c *gin.Context) {
	articleNumber_s := c.Param("articleNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
		return
	}
	err = models.RemoveBlogAritcle(articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	c.Redirect(http.StatusSeeOther, "/blog/1")
}
func BlogSearchRoute(c *gin.Context) {
	keyword := c.Query("search")
	articleNumber_s := c.Param("pageNumber")
	articleNumber, err := strconv.Atoi(articleNumber_s)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	searchedArticles, err := models.GetSearchedArticles(keyword, articleNumber)
	if err != nil {
		ErrorRediect(c, err.Error())
	}

	pageNumber_S := c.Param("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumber_S)
	if err != nil {
		ErrorRediect(c, "Bad request")
	}

	numberofBlogArticles, err := models.GetNumberOfSearchedArticles(keyword)
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	totalPages := numberofBlogArticles/6 + 1
	startpage := pageNumber - 2
	endpage := pageNumber + 2
	if startpage <= 1 {
		startpage = 1
	}
	if endpage >= totalPages {
		endpage = totalPages
	}
	var visiblePages []int
	for i := startpage; i <= endpage; i++ {
		visiblePages = append(visiblePages, i)
	}
	pagination := PageData{
		PageNumber:   pageNumber, // 실제 페이지 번호
		TotalPages:   totalPages, // 실제 전체 페이지 수
		VisiblePages: visiblePages,
	}
	isLoggedIn := c.GetBool("isLoggedIn")
	c.HTML(http.StatusOK, "blog-search.html", gin.H{
		"title":        "Blog Searched",
		"blogArticles": searchedArticles,
		"pagination":   pagination,
		"isLoggedIn":   isLoggedIn,
		"keyword":      keyword,
	})
}
