package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"regexp"
	"strings"
	"time"
)

//微信小程序
const (
	//code换session
	codeToSessionUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	//获得accessToken
	getAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

	//用户支付完成后，获取该用户的UnionId
	getPaidUnionidUrl = "https://api.weixin.qq.com/wxa/getpaidunionid?access_token=%s&openid=%s"
)

//WxAppletSessionData 微信小程序
type WxAppletSessionData struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WxAppletAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type WxAppletUserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

//AppletClient  Client
type AppletClient struct {
	AppId  string
	Secret string
}

//NewAppletClient
//传入appId和appSecret
func NewAppletClient(appId, secret string) *AppletClient {
	return &AppletClient{
		AppId:  appId,
		Secret: secret,
	}
}

//CodeToSession 根据code值换session
func (c *AppletClient) CodeToSession(code string) (sessionData *WxAppletSessionData, err error) {
	if code == "" {
		err = fmt.Errorf("[wechat]code为空")
		return
	}
	url := fmt.Sprintf(codeToSessionUrl, c.AppId, c.Secret, code)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	fmt.Println("resp:::", resp)
	sessionData = new(WxAppletSessionData)
	err = resp.ToJSON(&sessionData)
	if err != nil {
		return
	}
	if sessionData.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取code出错 %v", sessionData.ErrMsg)
		return
	}
	return
}

//GetAccessToken accessToken
func (c *AppletClient) GetAccessToken() (accessTokenData *WxAppletAccessToken, err error) {
	url := fmt.Sprintf(getAccessTokenUrl, c.AppId, c.Secret)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	accessTokenData = new(WxAppletAccessToken)
	err = resp.ToJSON(&accessTokenData)
	if err != nil {
		return
	}
	if accessTokenData.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取code出错 %v", accessTokenData.ErrMsg)
		return
	}
	return
}

func (c *AppletClient) Decrypt(sessionKey, encryptedData string, iv string) (*WxAppletUserInfo, error) {
	decrypted := new(WxAppletUserInfo)
	sessionKey = strings.Replace(strings.TrimSpace(sessionKey), " ", "+", -1)
	if len(sessionKey) != 24 {
		return decrypted, errors.New("sessionKey length is error")
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return decrypted, errors.New("DecodeString err")
	}
	iv = strings.Replace(strings.TrimSpace(iv), " ", "+", -1)
	if len(iv) != 24 {
		return decrypted, errors.New("iv length is error")
	}
	aesIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return decrypted, errors.New("DecodeString err")
	}
	encryptedData = strings.Replace(strings.TrimSpace(encryptedData), " ", "+", -1)
	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return decrypted, errors.New("DecodeString err")
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return decrypted, errors.New("NewCipher err")
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIv)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = PKCS7UnPadding(aesPlantText)

	//decrypted := make(map[string]interface{})

	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	aesPlantText = []byte(re.ReplaceAllString(string(aesPlantText), "$1"))
	err = json.Unmarshal(aesPlantText, &decrypted)
	if err != nil {
		return decrypted, errors.New("NewCipher err")
	}
	//if decrypted["watermark"].(map[string]interface{})["appid"] != c.AppId {
	//	return nil, errors.New("appId is not match")
	//}
	if decrypted.Watermark.AppID != c.AppId {
		return decrypted, errors.New("appId is not match")
	}
	return decrypted, nil
}

// PKCS7UnPadding return unpadding []Byte plantText
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}
