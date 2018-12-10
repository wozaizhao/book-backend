package server


import (
	"github.com/gin-gonic/gin"
	"wozaizhao.com/book/controllers/user"
	"wozaizhao.com/book/controllers/book"
)


func SetupRouter() *gin.Engine {
	r:= gin.Default()

	//小程序端接口
	wx := r.Group("/wx")
	{
	  //wx登录
	  wx.POST("/login",user.WxLogin)
	  //书籍列表
	  wx.GET("book",book.ListBooks)
	  
	  //书籍详情,包含目录
	  wx.GET("book/:id",book.GetBook)
	  //章节
	  wx.GET("book/:id/:content/:page")

		
	}

	//管理后台接口
	admin := r.Group("/admin")
	{
	  //后台登录	
	  admin.POST("login",user.Login)
	  //增加书籍
	  admin.POST("book",book.AddBook)
	  //增加目录
	  admin.POST("content",book.AddContent)
	  //增加页面
      admin.POST("page",book.AddPage)
	}

	return r
}