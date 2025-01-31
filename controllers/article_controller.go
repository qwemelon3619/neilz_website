package controllers

import (
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"neilz.space/web/models"
)

func BlogArticleRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "blog-article.html", gin.H{"title": "Blog Article"})
}
func BlogPostPageRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "blog-post.html", gin.H{"title": "Blog Posting"})
}
func imageResize(filepath string) error {
	file, err := os.Open(filepath)
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
	outFile, err := os.Create(filepath) // 덮어쓰기
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, croppedImg) // Encode the cropped image
	if err != nil {
		return err
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
	article.Content = c.PostForm("article-content")
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
	// c.JSON(http.StatusOK, article)
	Redirect(c, "/blog/1")
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
	totalPages := numberofBlogArticles % 6
	startpage := pageNumber - 2
	if startpage < 1 {
		startpage = 1
	}
	endpage := pageNumber + 2
	if endpage >= totalPages {
		startpage = totalPages - 2
		endpage = totalPages
	}
	var visiblePages []int
	for i := startpage; i <= endpage; i++ {
		visiblePages = append(visiblePages, i)
	}
	log.Print(visiblePages, pageNumber)
	pagination := PageData{
		PageNumber:   pageNumber, // 실제 페이지 번호
		TotalPages:   totalPages, // 실제 전체 페이지 수
		VisiblePages: visiblePages,
	}

	c.HTML(http.StatusOK, "blog-list.html", gin.H{
		"title":        "Blog List",
		"blogArticles": blogArticles,
		"pagination":   pagination,
	})
}
func AllBlogListJSON(c *gin.Context) {
	articles, err := models.GetAllBlogArticles()
	if err != nil {
		ErrorRediect(c, err.Error())
	}
	c.JSON(http.StatusOK, articles)
}
