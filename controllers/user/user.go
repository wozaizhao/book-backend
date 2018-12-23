package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wozaizhao.com/book/common"
	"wozaizhao.com/book/controllers/wechat"
	"wozaizhao.com/book/models"
	"wozaizhao.com/book/utils"
)

type WxLoginRequest struct {
	Code          string `json:"code" form:"code" binding:"required"`
	EncryptedData string `json:"encrypted_data" form:"encrypted_data" binding:"required"`
	Iv            string `json:"iv" form:"iv" binding:"required"`
}

type SubscribeReq struct {
	Status int `json:"status" form:"status" binding:"required"`
}

type Subscription struct {
	Self  bool
	Count int64
}

func WxLogin(c *gin.Context) {
	var wxl WxLoginRequest
	if err := c.Bind(&wxl); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": err.Error()})
		return
	} else {
		//小程序登录
		fmt.Println(wxl)

		wXBizDataCrypt, err := wechat.GetJsCode2Session(wxl.Code)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{"code": common.FAIL, "msg": "login fail"})
		}
		fmt.Println(wXBizDataCrypt)

		userinfo, err := wechat.WeDecryptData(wXBizDataCrypt, wxl.EncryptedData, wxl.Iv)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{"code": common.FAIL, "msg": "login fail"})
		}
		token := utils.Md5(wXBizDataCrypt.SessionKey)
		fmt.Println(token)

		//根据openid判断用户是否存在
		if models.UserExist(wXBizDataCrypt.Openid) {
			//存在则更新token
			models.UpdateUserToken(wXBizDataCrypt.Openid, wXBizDataCrypt.SessionKey, token)
			c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "msg": "relogin successfully!", "token": token, "userinfo": userinfo})
			return

		} else {
			//不存在则保存用户信息
			models.SaveUser(userinfo.OpenID, userinfo.NickName, userinfo.Gender, userinfo.City, userinfo.Province, userinfo.Country, userinfo.AvatarURL, userinfo.UnionID, wXBizDataCrypt.SessionKey, token)
			c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "msg": "login successfully!", "token": token, "userinfo": userinfo})
			return
		}

	}
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hi": "admin"})
}

func SubscribeBook(c *gin.Context) {
	bookid := c.Param("id")
	var subscribe SubscribeReq
	if c.ShouldBind(&subscribe) == nil {
		fmt.Println(subscribe.Status)
	}
	bookidint, err := strconv.Atoi(bookid)
	if err != nil {
		fmt.Println(err)
	}
	token := c.Request.Header["Token"][0]
	openid := models.Skey2OpenId(token)
	subscription := new(Subscription)
	if models.FavoriteExsit(openid, bookidint) {
		models.UpdateFavorite(openid, bookidint, subscribe.Status)
	} else {
		models.InsertFavorite(openid, bookidint, subscribe.Status)
	}
	if subscribe.Status == 1 {
		subscription.Self = true
	} else {
		subscription.Self = false
	}
	subscription.Count = models.Subscription(bookidint)
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "stared", "subscription": subscription})
}
