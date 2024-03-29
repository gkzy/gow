package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gkzy/gow/lib/util"
	"github.com/imroc/req"
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

//微信小程序用户基础信息
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

//微信小程序用户电话号码信息
type WxAppletPhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
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
	aesPlantText,err := util.AppletDecrypt(sessionKey,encryptedData,iv)
	if err != nil{
		return decrypted,err
	}
	err = json.Unmarshal(aesPlantText, &decrypted)
	if err != nil {
		return decrypted, err
	}
	if decrypted.Watermark.AppID != c.AppId {
		return decrypted, errors.New("appId is not match")
	}
	return decrypted, nil
}

func (c *AppletClient) DecryptPhoneInfo(sessionKey, encryptedData string, iv string) (*WxAppletPhoneInfo, error) {
	decrypted := new(WxAppletPhoneInfo)
	aesPlantText,err := util.AppletDecrypt(sessionKey,encryptedData,iv)
	if err != nil{
		return decrypted,err
	}
	err = json.Unmarshal(aesPlantText, &decrypted)
	if err != nil {
		return decrypted, errors.New("NewCipher err")
	}
	if decrypted.Watermark.AppID != c.AppId {
		return decrypted, errors.New("appId is not match")
	}
	return decrypted, nil
}

