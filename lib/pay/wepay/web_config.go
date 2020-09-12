package wepay

import "io/ioutil"

//WxConfig WxConfig
type WxConfig struct {
	appID     string //传入的appID
	mchID     string //分配的mchID
	apiKey    string //分配的apiKey
	serverIP  string //服务器IP
	certData  []byte //证书
	isSandbox bool   //是否沙箱
	notifyURL string //异步通知地址
	endMinute int    //订单有效分钟数
}

//NewWxConfig 一个新的配置信息
//也可以自己组装
func NewWxConfig(appID, mchID, apiKey, serverIP string, isSandbox bool) *WxConfig {
	return &WxConfig{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		serverIP:  serverIP,
		isSandbox: isSandbox,
	}
}

//SetCertData 设置证书
func (m *WxConfig) SetCertData(certPath string) (err error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}
	m.certData = certData
	return
}
