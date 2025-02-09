package controllers

import (
	"errors"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"neilz.space/web/models"
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
			"title":       "Blog Posting",
			"blogArticle": article,
			"content":     content,
			"isLoggedIn":  isLoggedIn,
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
func imageResize(filepath_ string) error {
	file, err := os.Open(filepath_)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 1. 크기 조정
	newImg := resize(img, 750, 350)

	// 2. 자르기 (필요한 경우)
	rect := image.Rect(0, 0, 750, 350)
	croppedImg := newImg.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	// 3. 파일 저장 (원래 파일에 덮어쓰기)
	outFile, err := os.Create(filepath_) // 덮어쓰기
	if err != nil {
		return err
	}
	defer outFile.Close()
	ext := filepath.Ext(filepath_)
	log.Print(ext)
	if ext == ".png" {
		err = png.Encode(outFile, croppedImg) // Encode the cropped image
		if err != nil {
			return err
		}
	} else if ext == ".jpeg" || ext == ".jpg" {
		err = jpeg.Encode(outFile, croppedImg, nil) // Encode the cropped image
		if err != nil {
			return err
		}
	} else {
		return errors.New("file does not support")
	}

	return nil
}
func resize(img image.Image, width, height int) image.Image {
	// 새로운 이미지 생성
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// 이미지 비율에 맞춰서 조정
	srcBounds := img.Bounds()
	srcAspect := float64(srcBounds.Dx()) / float64(srcBounds.Dy())
	dstAspect := float64(width) / float64(height)

	var drawRect image.Rectangle
	if srcAspect > dstAspect {
		newHeight := int(float64(width) / srcAspect)
		yOffset := (height - newHeight) / 2
		drawRect = image.Rect(0, yOffset, width, yOffset+newHeight)
	} else {
		newWidth := int(float64(height) * srcAspect)
		xOffset := (width - newWidth) / 2
		drawRect = image.Rect(xOffset, 0, xOffset+newWidth, height)
	}

	// 이미지 그리기 (크기 조정)
	draw.Draw(newImg, drawRect, img, srcBounds.Min, draw.Src)

	return newImg
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

	err = imageResize(filePath) //crop image
	if err != nil {
		ErrorRediect(c, err.Error())
		return
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
