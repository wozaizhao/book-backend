package user

import (
	"net/http"

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
		common.Log("Request of wxl:", wxl)

		wXBizDataCrypt, err := wechat.GetJsCode2Session(wxl.Code)
		if err != nil {
			common.Log("GetJsCode2Session Error:", err)
			c.JSON(http.StatusOK, gin.H{"code": common.FAIL, "msg": "login fail"})
			return
		}
		common.Log("wXBizDataCrypt", wXBizDataCrypt)

		userinfo, err := wechat.WeDecryptData(wXBizDataCrypt, wxl.EncryptedData, wxl.Iv)
		if err != nil {
			common.Log("WeDecryptData Error:", err)
			c.JSON(http.StatusOK, gin.H{"code": common.FAIL, "msg": "login fail"})
			return
		}
		common.Log("userinfo", userinfo)
		token := utils.Md5(wXBizDataCrypt.SessionKey)
		common.Log("token", token)

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
	id := c.Param("id")
	common.Log("SubscribeBook id:", id)
	var subscribe SubscribeReq
	if err := c.ShouldBind(&subscribe); err != nil {
		common.Log("ShouldBind Error:", err)
	}
	common.Log("SubscribeBook status:", subscribe.Status)
	bookid := common.String2int(id)
	token := c.Request.Header["Token"]
	if token == nil {
		c.JSON(http.StatusOK, gin.H{"code": common.FAIL, "message": "token missing"})
		return
	}
	common.Log("SubscribeBook token:", token)
	openid := models.Skey2OpenId(token[0])
	if openid == "" {
		c.JSON(http.StatusOK, gin.H{"code": common.TOKENEXPIRED, "message": "token expired"})
		return
	}
	subscription := new(Subscription)
	if models.FavoriteExsit(openid, bookid) {
		models.UpdateFavorite(openid, bookid, subscribe.Status)
	} else {
		models.InsertFavorite(openid, bookid, subscribe.Status)
	}
	if subscribe.Status == common.SUBSCRIBE {
		subscription.Self = true
	} else {
		subscription.Self = false
	}
	subscription.Count = models.Subscription(bookid)
	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESS, "message": "stared", "subscription": subscription})
}
