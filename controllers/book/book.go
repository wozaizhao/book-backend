package book

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"wozaizhao.com/book/common"
	"wozaizhao.com/book/models"
)

type Book struct {
	Priority int    `form:"priority" json:"priority" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Cate     string `form:"type" json:"cate" binding:"required"`
	Cover    string `form:"cover" json:"cover" binding:"required"`
	Slogan   string `form:"slogan" json:"slogan" binding:"required"`
	Bg       string `form:"bg" json:"bg"`
	Color    string `form:"color" json:"color"`
	Tag      string `form:"tag" json:"tag"`
	Intro    string `form:"intro" json:"intro" binding:"required"`
	Path     string `form:"path" json:"path" binding:"required"`
	Url      string `form:"url" json:"url"`
}

type Content struct {
	BookId int    `form:"bookid" json:"bookid" binding:"required"`
	Sn     int    `form:"sn" json:"sn" binding:"required"`
	Title  string `form:"title" json:"title" binding:"required"`
	Pages  []models.Page
}

type Page struct {
	BookId    int `form:"bookid" json:"bookid" binding:"required"`
	ContentId int `form:"contentid" json:"contentid" binding:"required"`
	// Sn int         `form:"sn" json:"sn" binding:"required"`
	Title string `form:"title" json:"title" binding:"required"`
	MdUrl string `form:"mdurl" json:"mdurl" binding:"required"`
}

type Pagination struct {
	BookId    int
	ContentId int
	Sn        int
}

type Pages struct {
	BookId    int        `form:"bookid" json:"bookid" binding:"required"`
	ContentId int        `form:"contentid" json:"contentid" binding:"required"`
	Path      string     `form:"path" json:"path" binding:"required"`
	PageArray []PageItem `form:"pages" json:"pages" binding:"required"`
}

type PageItem struct {
	Title string `form:"title" json:"title"`
	MdUrl string `form:"mdurl" json:"mdurl"`
}

type Subscription struct {
	Self  bool
	Count int64
}

func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		common.Log("AddBook ShouldBindJSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": common.FAIL, "message": "AddBook BadRequest"})
		return
	}
	b := models.InsertBook(book.Priority, book.Name, book.Cate, book.Cover, book.Slogan, book.Bg, book.Color, book.Tag, book.Intro, book.Path, book.Url)
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "AddBook Successfully", "book": b})
}

func ListBooks(c *gin.Context) {
	token := c.Request.Header["Token"]
	var selfbooks []models.Book

	var books = models.GetBooks()
	// fmt.Println(books)
	if token != nil {
		openid := models.Skey2OpenId(token[0])
		fmt.Println(openid == "")
		if openid == "" {
			c.JSON(http.StatusOK, gin.H{"books": books, "selfbooks": nil})
			return
		}
		// fmt.Println(openid)
		//用户是否已订阅本书
		selfbooks = models.GetSelfBooks(openid)
		// fmt.Println(selfbooks)
	}
	// fmt.Println(books)
	c.JSON(http.StatusOK, gin.H{"books": books, "selfbooks": selfbooks})

}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	bookid := common.String2int(id)

	subscription := new(Subscription)

	token := c.Request.Header["Token"]
	common.Log("token", token)
	if token != nil {
		openid := models.Skey2OpenId(token[0])
		//用户是否已订阅本书
		if openid == "" {
			subscription.Self = false
		} else {
			subscription.Self = models.HasUserSubscribe(openid, bookid)
		}
	}

	//本书总订阅人数
	subscription.Count = models.Subscription(bookid)
	common.Log("GetBook subscription:", subscription)
	book := models.GetBook(bookid)
	contents := models.GetContents(bookid)
	theContents := make([]Content, len(contents))
	for i := 0; i < len(contents); i++ {
		theContents[i].BookId = contents[i].BookId
		theContents[i].Sn = contents[i].Sn
		theContents[i].Title = contents[i].Title
		theContents[i].Pages = models.GetPages(contents[i].Id)
		// fmt.Println(pages)
		// theContents[i].Pages = pages
	}
	c.JSON(http.StatusOK, gin.H{"book": book, "content": theContents, "subscription": subscription})
}

func AddContent(c *gin.Context) {
	var content Content
	if err := c.ShouldBindJSON(&content); err != nil {
		common.Log("AddContent ShouldBindJSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": common.FAIL, "message": "AddContent BadRequest"})
		return
	}
	ct := models.InsertContent(content.BookId, content.Sn, content.Title)
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "AddContent Successfully", "content": ct})
}

func AddPage(c *gin.Context) {
	var page Page
	if err := c.ShouldBindJSON(&page); err != nil {
		common.Log("AddPage ShouldBindJSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": common.FAIL, "message": "AddPage BadRequest"})
		return
	}
	p := models.InsertPage(page.BookId, page.ContentId, page.Title, page.MdUrl)
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "AddPage Successfully", "page": p})
}

func AddPages(c *gin.Context) {
	var pages Pages
	if err := c.ShouldBindJSON(&pages); err != nil {
		common.Log("AddPages ShouldBindJSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": common.FAIL, "message": "AddPages BadRequest"})
		return
	}
	bookid := pages.BookId
	contentid := pages.ContentId
	pagearray := pages.PageArray
	path := pages.Path
	for i := 0; i < len(pagearray); i++ {
		// fmt.Println(bookid,contentid, pagearray[i].Title , path + pagearray[i].MdUrl)
		models.InsertPage(bookid, contentid, pagearray[i].Title, path+pagearray[i].MdUrl)
	}
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "AddPages Successfully", "count": len(pagearray)})
}

func GetPage(c *gin.Context) {
	id := c.Param("id")
	contentid := c.Param("content")
	pageid := c.Param("page")
	page := models.GetPage(id, contentid, pageid)
	next := hasNextPage(id, contentid, pageid)
	prev := hasPrevPage(id, contentid, pageid)
	c.JSON(http.StatusOK, gin.H{"page": page, "next": next, "prev": prev})
}

func GetPageById(c *gin.Context) {
	id := c.Param("id")
	page := models.GetPageById(id)
	c.JSON(http.StatusOK, gin.H{"page": page})
}

func hasNextPage(bookid string, contentid string, sn string) (p *Pagination) {
	book_id := common.String2int(bookid)
	content_id := common.String2int(contentid)
	page_sn := common.String2int(sn)

	// if errb != nil || errc != nil || errp != nil {
	// 	fmt.Println("somethis wrong")
	// }

	page_sn = page_sn + 1

	has := models.PageExist(book_id, content_id, page_sn)

	if has {
		p = new(Pagination)
		p.BookId = book_id
		p.ContentId = content_id
		p.Sn = page_sn
	}
	return

}

func hasPrevPage(bookid string, contentid string, sn string) (p *Pagination) {
	book_id := common.String2int(bookid)
	content_id := common.String2int(contentid)
	page_sn := common.String2int(sn)

	page_sn = page_sn - 1

	has := models.PageExist(book_id, content_id, page_sn)

	if has {
		p = new(Pagination)
		p.BookId = book_id
		p.ContentId = content_id
		p.Sn = page_sn

	}
	return
}

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"search": "result"})
}
