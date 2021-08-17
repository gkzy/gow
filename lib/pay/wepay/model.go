package wepay

import "encoding/xml"

//AppPayResp App支付时，返回的结构体
//包括prepayID和Sign等其他信息
type AppPayResp struct {
	AppID     string `json:"appid"`     //appid
	PartnerID string `json:"partnerid"` //商户ID
	PrepayID  string `json:"prepayid"`  //prepayid
	Package   string `json:"package"`   //package
	NonceStr  string `json:"noncestr"`  //随机字串
	Timestamp string `json:"timestamp"` //时间
	Sign      string `json:"sign"`      //签名
}

//NotifyRet 异步通知的返回值
//返回
type NotifyRet struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

//AppletPayResp 微信小程序支付时，返回的结构体
//包括prepayID和Sign等其他信息
type AppletPayResp struct {
	Timestamp string `json:"timeStamp"` //时间
	NonceStr  string `json:"nonceStr"`  //随机字串
	Package   string `json:"package"`   //package
	SignType  string `json:"signType"`
	Sign      string `json:"paySign"` //签名
}

// RefundNotifyReqInfo 退款异步通知加密字段
type RefundNotifyReqInfo struct {
	TransactionId       string `json:"transaction_id"`
	OutTradeNo          string `json:"out_trade_no"`
	RefundId            string `json:"refund_id"`
	OutRefundNo         string `json:"out_refund_no"`
	TotalFee            int64  `json:"total_fee"`
	SettlementTotalFee  int64  `json:"settlement_total_fee"`  //应结订单金额
	RefundFee           int64  `json:"refund_fee"`            //申请退款金额
	SettlementRefundFee int64  `json:"settlement_refund_fee"` //退款金额
	RefundStatus        string `json:"refund_status"`         //退款状态 SUCCESS-退款成功 CHANGE-退款异常 REFUNDCLOSE—退款关闭
	RefundRecvAccout    string `json:"refund_recv_accout"`    //退款入账账户
	RefundAccount       string `json:"refund_account"`        //退款资金来源 REFUND_SOURCE_RECHARGE_FUNDS 可用余额退款/基本账户  REFUND_SOURCE_UNSETTLED_FUNDS 未结算资金退款
	RefundRequestSource string `json:"refund_request_source"` //退款发起来源 API接口  VENDOR_PLATFORM商户平台
}
