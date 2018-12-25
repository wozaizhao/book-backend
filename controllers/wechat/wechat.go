package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"wozaizhao.com/book/common"
	"wozaizhao.com/book/utils"
)

type (

	// WxUserInfo 微信用户资料
	// WxUserInfo struct {
	// 	OpenID     string `json:"openid,omitempty"`     // 授权用户唯一标识
	// 	NickName   string `json:"nickname,omitempty"`   // 普通用户昵称
	// 	Sex        uint32 `json:"sex,omitempty"`        // 普通用户性别，1为男性，2为女性
	// 	Province   string `json:"province,omitempty"`   // 普通用户个人资料填写的省份
	// 	City       string `json:"city,omitempty"`       // 普通用户个人资料填写的城市
	// 	Country    string `json:"country,omitempty"`    // 国家，如中国为CN
	// 	HeadImgURL string `json:"headimgurl,omitempty"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	// 	//Privilege  string `json:"privilege"`
	// 	Privilege []string `json:"privilege,omitempty"` // 用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	// 	UnionID   string   `json:"unionid,omitempty"`   // 普通用户的标识，对当前开发者帐号唯一
	// 	ErrCode   uint     `json:"errcode,omitempty"`
	// 	ErrMsg    string   `json:"errmsg,omitempty"`
	// }

	// WechatEncryptedData 小程序解密后结构
	WechatEncryptedData struct {
		OpenID    string          `json:"openId"`
		NickName  string          `json:"nickName"`
		Gender    int             `json:"gender"` //性别，0-未知，1-男，2-女
		City      string          `json:"city"`
		Province  string          `json:"province"`
		Country   string          `json:"country"`
		AvatarURL string          `json:"avatarUrl"`
		UnionID   string          `json:"unionId"`
		WaterMark WechatWaterMark `json:"watermark"` //水印
	}

	// WechatWaterMark 加密验证信息
	WechatWaterMark struct {
		Appid     string `json:"appid"`
		Timestamp uint64 `json:"timestamp"`
	}

	// WXBizDataCrypt 小程序解密密钥信息
	WXBizDataCrypt struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
)

//获取appid等
func getAppid() (appid, appSecret string) {
	err := godotenv.Load()
	if err != nil {
		common.Log("getAppid Error:", err)
	}

	appid = os.Getenv("APPID")
	appSecret = os.Getenv("APPSECRET")

	return
}

// WxLogin 微信小程序登录 直接登录获取用户信息
// func WxLogin(code, encryptedData, iv string) (wxUserInfo *WechatEncryptedData, err error) {
// 	wXBizDataCrypt, err := GetJsCode2Session(code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(wXBizDataCrypt)
// 	return WeDecryptData(wXBizDataCrypt, encryptedData, iv)
// }

// GetJsCode2Session 获取
func GetJsCode2Session(code string) (wXBizDataCrypt *WXBizDataCrypt, err error) {

	if code == "" {
		return wXBizDataCrypt, errors.New("GetJsCode2Session error: code is null")
	}

	appid, appSecret := getAppid()

	params := url.Values{}
	params.Add("appid", appid)
	params.Add("secret", appSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	// params := url.Values{
	// 	"js_code":    []string{code},
	// 	"grant_type": []string{"authorization_code"},
	// }

	body, err := utils.NewRequest("GET", common.JsCode2SessionURL, []byte(params.Encode()))
	if err != nil {
		return wXBizDataCrypt, err
	}
	err = json.Unmarshal(body, &wXBizDataCrypt)
	if err != nil {
		return wXBizDataCrypt, err
	}

	if wXBizDataCrypt.ErrMsg != "" {
		return wXBizDataCrypt, errors.New(wXBizDataCrypt.ErrMsg)
	}

	return
}

// WeDecryptData 微信小程序登录数据解密
func WeDecryptData(wXBizDataCrypt *WXBizDataCrypt, encryptedData, iv string) (wechatEncryptedData *WechatEncryptedData, err error) {

	if len(wXBizDataCrypt.SessionKey) != 24 {
		return wechatEncryptedData, errors.New("encodingAesKey illegal")
	}

	aesKey, err := base64.StdEncoding.DecodeString(wXBizDataCrypt.SessionKey)
	if err != nil {
		return wechatEncryptedData, err
	}
	aesCipher, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return wechatEncryptedData, err
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return wechatEncryptedData, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		errMsg := fmt.Sprintf("aes new cipher error: %#v", err)
		return wechatEncryptedData, errors.New(errMsg)
	}

	c := cipher.NewCBCDecrypter(block, aesIV)
	resBytes := make([]byte, len(aesCipher))
	c.CryptBlocks(resBytes, aesCipher)
	resBytes = utils.PKCS7UnPadding(resBytes)

	//解密后的byte数组数据做json解析
	wechatEncryptedData = &WechatEncryptedData{}
	err = json.Unmarshal(resBytes, &wechatEncryptedData)
	if err != nil {
		errMsg := fmt.Sprintf("json unmarshal data error: %#v", err)
		return wechatEncryptedData, errors.New(errMsg)
	}

	return
}
