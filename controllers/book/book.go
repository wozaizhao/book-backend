package book

import (
	"github.com/gin-gonic/gin"
	"wozaizhao.com/book/models"
	"net/http"
	"fmt"
	"strconv"

)

type Book struct {
	Priority int  `form:"priority" json:"priority" binding:"required"`
	Name string   `form:"name" json:"name" binding:"required"`
	Cate string   `form:"type" json:"cate" binding:"required"`
	Cover string  `form:"cover" json:"cover" binding:"required"`
	Slogan string `form:"slogan" json:"slogan" binding:"required"`
	Bg string     `form:"bg" json:"bg"`
	Color string  `form:"color" json:"color"`
	Tag string    `form:"tag" json:"tag"`
	Intro string  `form:"intro" json:"intro" binding:"required"`
	Path string   `form:"path" json:"path" binding:"required"`
}

type Content struct {
	BookId int     `form:"bookid" json:"bookid" binding:"required"`
	Sn int         `form:"sn" json:"sn" binding:"required"`
	Title string   `form:"title" json:"title" binding:"required"`
	Pages []models.Page
}

type Page struct {
	BookId int     `form:"bookid" json:"bookid" binding:"required"`
	ContentId int  `form:"contentid" json:"contentid" binding:"required"`
	Sn int         `form:"sn" json:"sn" binding:"required"`
	Title string   `form:"title" json:"title" binding:"required"`
	MdUrl string   `form:"mdurl" json:"mdurl" binding:"required"`
}

func AddBook(c *gin.Context){
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	models.InsertBook(book.Priority,book.Name,book.Cate,book.Cover,book.Slogan,book.Bg,book.Color,book.Tag,book.Intro,book.Path)
	c.JSON(http.StatusOK,gin.H{"addbook":"ok"})
}

func ListBooks(c *gin.Context) {
	var books = models.GetBooks()
	fmt.Println(books)
	if books != nil {
		c.JSON(http.StatusOK,gin.H{"books":books})
	} else {
		c.JSON(http.StatusOK,gin.H{"books":nil})
	}
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	idint, err :=  strconv.Atoi(id)
	if err!= nil {
		fmt.Println(err)
	}
	book := models.GetBook(idint)
	contents := models.GetContents(idint)
	theContents := make([]Content, len(contents))
	for i:=0; i < len(contents); i++ {
		theContents[i].BookId = contents[i].BookId
		theContents[i].Sn = contents[i].Sn
		theContents[i].Title = contents[i].Title
		pages := models.GetPages(contents[i].Id)
		// fmt.Println(pages)
		theContents[i].Pages = pages
    }

	if book != nil {
		c.JSON(http.StatusOK,gin.H{"book":book,"content":theContents})
	} else {
		c.JSON(http.StatusOK,gin.H{"book":nil})
	}
}

func AddContent(c *gin.Context){
	var content Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	models.InsertContent(content.BookId,content.Sn,content.Title)
	c.JSON(http.StatusOK,gin.H{"addcontent":"ok"})
}

func AddPage(c *gin.Context){
	var page Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	models.InsertPage(page.BookId,page.ContentId ,page.Sn , page.Title , page.MdUrl)
	c.JSON(http.StatusOK,gin.H{"addpage":"ok"})
}