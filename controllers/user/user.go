package user

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type WxLoginRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
	EncryptedData string `json:"encrypted_data" form:"encrypted_data" binding:"required"`
	Iv string `json:"iv" form:"iv" binding:"required"`
}

func WxLogin(c *gin.Context){
	var wxl WxLoginRequest
	if err := c.Bind(&wxl); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": err.Error()})
		return
	} else {
		fmt.Println("参数：",wxl.Code)
		fmt.Println("参数：",wxl.EncryptedData)
		fmt.Println("参数：",wxl.Iv)
		
		c.JSON(http.StatusOK, gin.H{"code":0,"msg":"success", "token": "token"})

	}
}

func Login(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{"hi":"admin"})
}