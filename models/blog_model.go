package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

//	type Tag struct {
//		tag_id   int
//		tag_name string
//	}
//
//	type Post_Tags struct {
//		post_id int `gorm:"forignKey:BlogPost;constraint:OnDelete:CASCADE"`
//		tag_id  int `gorm:"forignKey:Tag"`
//	}
type Author struct {
	Name            string
	Job             string
	ProfileImageURL string
}
type BlogArticleInput struct {
	Title           string `gorm:"not null;"`
	SubTitle        string
	Content         string `gorm:"not null;"`
	HeaderImageName string
}
type BlogArticle struct {
	gorm.Model
	Author
	Title           string
	Subtitle        string
	Content         string `gorm:"type:text"`
	HeaderImageName string
}

func SaveBlogArticle(input BlogArticleInput) error {
	article := BlogArticle{}
	article.Title = input.Title
	article.Subtitle = input.SubTitle
	article.Content = input.Content
	article.Author.Name = "Neil"
	article.Author.Job = "Backend Developer"
	article.HeaderImageName = input.HeaderImageName

	_, err := article.SaveArticleInDB()

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return errors.New("DataBase Error: could not save the data")
	}
	return nil
}
func EditBlogArticle(input BlogArticleInput, articleNumber int) error {
	article := BlogArticle{}
	article.ID = uint(articleNumber)
	article.Title = input.Title
	article.Subtitle = input.SubTitle
	article.Content = input.Content
	article.Author.Name = "Neil"
	article.Author.Job = "Backend Developer"
	if input.HeaderImageName != "" {
		article.HeaderImageName = input.HeaderImageName
	}

	_, err := article.EditArticleInDB()

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return errors.New("DataBase Error: could not save the data")
	}
	return nil
}
func RemoveBlogAritcle(articleNumber int) error {
	article := BlogArticle{}
	article.ID = uint(articleNumber)
	_, err := article.RemoveArticleInDB()
	if err != nil {
		return errors.New("DataBase Error: could not save the data")
	}
	return nil
}

func (a *BlogArticle) SaveArticleInDB() (*BlogArticle, error) {
	err := DB.Create(&a).Error
	if err != nil {
		return &BlogArticle{}, err
	}
	return a, nil
}
func (a *BlogArticle) EditArticleInDB() (*BlogArticle, error) {
	err := DB.Model(&a).Update(a).Error
	if err != nil {
		return &BlogArticle{}, err
	}
	return a, nil
}
func (a *BlogArticle) RemoveArticleInDB() (*BlogArticle, error) {
	err := DB.Delete(&a).Error
	if err != nil {
		return &BlogArticle{}, err
	}
	return a, nil
}
func GetAllBlogArticles() ([]BlogArticle, error) {
	var articles []BlogArticle
	result := DB.Model([]BlogArticle{}).Order("ID desc").Find(&articles)
	if result.Error != nil {
		return []BlogArticle{}, errors.New("DB error")
	}
	return articles, nil
}
func GetPageBlogArticles(pageNumber int) ([]BlogArticle, error) {
	var articles []BlogArticle
	maxNumberOfArticles := 6
	result := DB.Model([]BlogArticle{}).Order("ID desc").Find(&articles)
	if result.Error != nil {
		return []BlogArticle{}, errors.New("DB error")
	}
	var slicedArticles []BlogArticle
	if len(articles) > maxNumberOfArticles {
		start := (pageNumber - 1) * maxNumberOfArticles
		end := pageNumber * maxNumberOfArticles
		if pageNumber*maxNumberOfArticles >= len(articles) {
			end = len(articles)
		}
		slicedArticles = articles[start:end]
	} else {
		slicedArticles = articles
	}
	return slicedArticles, nil
}
func GetBlogArticle(id int) (BlogArticle, error) {
	var article BlogArticle
	result := DB.First(&article, id)
	if result.Error != nil {
		return BlogArticle{}, errors.New("DB error")
	}
	return article, nil
}
func GetNumberOfBlogArticles() (int, error) {
	var numberOfArticle int
	DB.Model([]BlogArticle{}).Count(&numberOfArticle)
	return numberOfArticle, nil
}
