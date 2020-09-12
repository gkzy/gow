package sms

// ISMS 接口
type ISMS interface {
	// SendVerifyCode 发送验证码短信
	//	phones为多个时，用,分割
	SendVerifyCode(phones, code string, templateId int) (bool, error)

	// SendNotice 发送通知短信
	SendNotice(phones, code string, templateId int) (bool, error)

	// SendSell 发送营销短信
	SendSell(phones, code string, templateId int) (bool, error)

	//StatusCallBack 状态回调
	StatusCallBack()
}
