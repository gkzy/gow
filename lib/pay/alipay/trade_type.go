package alipay

import "encoding/json"

//TradePay TradePay
type TradePay struct {
	NotifyURL string `json:"-"`
	ReturnURL string `json:"-"`

	// biz content，这四个参数是必须的
	Subject     string `json:"subject"`      // 订单标题
	OutTradeNo  string `json:"out_trade_no"` // 商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复
	TotalAmount string `json:"total_amount"` // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	ProductCode string `json:"product_code"` // 销售产品码，与支付宝签约的产品码名称。 注：目前仅支持FAST_INSTANT_TRADE_PAY

	Body               string `json:"body,omitempty"`                 // 订单描述
	BusinessParams     string `json:"business_params,omitempty"`      // 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	DisablePayChannels string `json:"disable_pay_channels,omitempty"` // 禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 注，与enable_pay_channels互斥
	EnablePayChannels  string `json:"enable_pay_channels,omitempty"`  // 可用渠道，用户只能在指定渠道范围内支付  当有多个渠道时用“,”分隔 注，与disable_pay_channels互斥
	//ExtUserInfo        string `json:"ext_user_info,omitempty"`        // 外部指定买家
	ExtendParams     string `json:"extend_params,omitempty"`     // 业务扩展参数，详见下面的“业务扩展参数说明”
	GoodsType        string `json:"goods_type,omitempty"`        // 商品主类型：0—虚拟类商品，1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	InvoiceInfo      string `json:"invoice_info,omitempty"`      // 开票信息
	PassbackParams   string `json:"passback_params,omitempty"`   // 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PromoParams      string `json:"promo_params,omitempty"`      // 优惠参数 注：仅与支付宝协商后可用
	RoyaltyInfo      string `json:"royalty_info,omitempty"`      // 描述分账信息，json格式，详见分账参数说明
	SellerID         string `json:"seller_id,omitempty"`         // 收款支付宝用户ID。 如果该值为空，则默认为商户签约账号对应的支付宝用户ID
	SettleInfo       string `json:"settle_info,omitempty"`       // 描述结算信息，json格式，详见结算参数说明
	SpecifiedChannel string `json:"specified_channel,omitempty"` // 指定渠道，目前仅支持传入pcredit  若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。  注：该参数不可与花呗分期参数同时传入
	StoreID          string `json:"store_id,omitempty"`          // 商户门店编号。该参数用于请求参数中以区分各门店，非必传项。
	SubMerchant      string `json:"sub_merchant,omitempty"`      // 间连受理商户信息体，当前只对特殊银行机构特定场景下使用此字段
	TimeoutExpress   string `json:"timeout_express,omitempty"`   // 该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
}

// TradePagePay  TradePagePay
// DOC https://docs.open.alipay.com/api_1/alipay.trade.app.pay
type TradePagePay struct {
	TradePay
	AuthToken   string `json:"auth_token,omitempty"`   // 针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
	GoodsDetail string `json:"goods_detail,omitempty"` // 订单包含的商品列表信息，Json格式，详见商品明细说明
	QRPayMode   string `json:"qr_pay_mode,omitempty"`  // PC扫码支付的方式，支持前置模式和跳转模式。
	QRCodeWidth string `json:"qrcode_width,omitempty"` // 商户自定义二维码宽度 注：qr_pay_mode=4时该参数生效
}

//APIName APIName
func (t TradePagePay) APIName() string {
	return "alipay.trade.page.pay"
}

//Params Params
func (t TradePagePay) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	m["return_url"] = t.ReturnURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradePagePay) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradePagePay) ExtJSONParamValue() string {
	var bytes, err = json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(bytes)
}

//=======================================================================================================

const (
	K_TRADE_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY" //（交易创建，等待买家付款）
	K_TRADE_STATUS_TRADE_CLOSED   = "TRADE_CLOSED"   //（未付款交易超时关闭，或支付完成后全额退款）
	K_TRADE_STATUS_TRADE_SUCCESS  = "TRADE_SUCCESS"  //（交易支付成功）
	K_TRADE_STATUS_TRADE_FINISHED = "TRADE_FINISHED" //（交易结束，不可退款）
)

// TradeQuery 订单查询
// https://docs.open.alipay.com/api_1/alipay.trade.query/
type TradeQuery struct {
	AppAuthToken string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 订单支付时传入的商户订单号, 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 支付宝交易号
}

//APIName APIName
func (t TradeQuery) APIName() string {
	return "alipay.trade.query"
}

//Params Params
func (t TradeQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeQuery) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeQuery) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeQueryResponse 订单查询返回信息
type TradeQueryResponse struct {
	AliPayTradeQuery struct {
		Code             string `json:"code"`
		Msg              string `json:"msg"`
		SubCode          string `json:"sub_code"`
		SubMsg           string `json:"sub_msg"`
		AuthTradePayMode string `json:"auth_trade_pay_mode"` // 预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
		BuyerLogonID     string `json:"buyer_logon_id"`      // 买家支付宝账号
		BuyerPayAmount   string `json:"buyer_pay_amount"`    // 买家实付金额，单位为元，两位小数。
		BuyerUserID      string `json:"buyer_user_id"`       // 买家在支付宝的用户id
		BuyerUserType    string `json:"buyer_user_type"`     // 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
		InvoiceAmount    string `json:"invoice_amount"`      // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		OutTradeNo       string `json:"out_trade_no"`        // 商家订单号
		PointAmount      string `json:"point_amount"`        // 积分支付的金额，单位为元，两位小数。
		ReceiptAmount    string `json:"receipt_amount"`      // 实收金额，单位为元，两位小数
		SendPayDate      string `json:"send_pay_date"`       // 本次交易打款给卖家的时间
		TotalAmount      string `json:"total_amount"`        // 交易的订单金额
		TradeNo          string `json:"trade_no"`            // 支付宝交易号
		TradeStatus      string `json:"trade_status"`        // 交易状态

		DiscountAmount      string          `json:"discount_amount"`               // 平台优惠金额
		FundBillList        []FundBill      `json:"fund_bill_list,omitempty"`      // 交易支付使用的资金渠道
		MdiscountAmount     string          `json:"mdiscount_amount"`              // 商家优惠金额
		PayAmount           string          `json:"pay_amount"`                    // 支付币种订单金额
		PayCurrency         string          `json:"pay_currency"`                  // 订单支付币种
		SettleAmount        string          `json:"settle_amount"`                 // 结算币种订单金额
		SettleCurrency      string          `json:"settle_currency"`               // 订单结算币种
		SettleTransRate     string          `json:"settle_trans_rate"`             // 结算币种兑换标价币种汇率
		StoreID             string          `json:"store_id"`                      // 商户门店编号
		StoreName           string          `json:"store_name"`                    // 请求交易支付中的商户店铺的名称
		TerminalID          string          `json:"terminal_id"`                   // 商户机具终端编号
		TransCurrency       string          `json:"trans_currency"`                // 标价币种
		TransPayRate        string          `json:"trans_pay_rate"`                // 标价币种兑换支付币种汇率
		DiscountGoodsDetail string          `json:"discount_goods_detail"`         // 本次交易支付所使用的单品券优惠的商品优惠信息
		IndustrySepcDetail  string          `json:"industry_sepc_detail"`          // 行业特殊信息（例如在医保卡支付业务中，向用户返回医疗信息）。
		VoucherDetailList   []VoucherDetail `json:"voucher_detail_list,omitempty"` // 本交易支付时使用的所有优惠券信息
	} `json:"alipay_trade_query_response"`
	Sign string `json:"sign"`
}

//FundBill FundBill
type FundBill struct {
	FundChannel string  `json:"fund_channel"`       // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string  `json:"amount"`             // 该支付工具类型所使用的金额
	RealAmount  float64 `json:"real_amount,string"` // 渠道实际付款金额
}

//VoucherDetail VoucherDetail
type VoucherDetail struct {
	ID                 string `json:"id"`                  // 券id
	Name               string `json:"name"`                // 券名称
	Type               string `json:"type"`                // 当前有三种类型： ALIPAY_FIX_VOUCHER - 全场代金券, ALIPAY_DISCOUNT_VOUCHER - 折扣券, ALIPAY_ITEM_VOUCHER - 单品优惠
	Amount             string `json:"amount"`              // 优惠券面额，它应该会等于商家出资加上其他出资方出资
	MerchantContribute string `json:"merchant_contribute"` // 商家出资（特指发起交易的商家出资金额）
	OtherContribute    string `json:"other_contribute"`    // 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资
	Memo               string `json:"memo"`                // 优惠券备注信息
}

//IsSuccess IsSuccess
func (t *TradeQueryResponse) IsSuccess() bool {
	if t.AliPayTradeQuery.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//=======================================================================================================

// TradeClose 统一收单交易关闭接口
// https://docs.open.alipay.com/api_1/alipay.trade.close/
type TradeClose struct {
	AppAuthToken string `json:"-"`                      // 可选
	NotifyURL    string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	OperatorID   string `json:"operator_id,omitempty"`  // 可选
}

// APIName APIName
func (t TradeClose) APIName() string {
	return "alipay.trade.close"
}

//Params Params
func (t TradeClose) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeClose) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeClose) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeCloseResponse TradeCloseResponse
type TradeCloseResponse struct {
	AliPayTradeClose struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		SubCode    string `json:"sub_code"`
		SubMsg     string `json:"sub_msg"`
		OutTradeNo string `json:"out_trade_no"`
		TradeNo    string `json:"trade_no"`
	} `json:"alipay_trade_close_response"`
	Sign string `json:"sign"`
}

//=======================================================================================================

// TradeRefund  统一收单交易退款接口
// https://docs.open.alipay.com/api_1/alipay.trade.refund/
type TradeRefund struct {
	AppAuthToken string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	RefundAmount string `json:"refund_amount"`          // 必须 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundReason string `json:"refund_reason"`          // 可选 退款的原因说明
	OutRequestNo string `json:"out_request_no"`         // 可选 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorID   string `json:"operator_id"`            // 可选 商户的操作员编号
	StoreID      string `json:"store_id"`               // 可选 商户的门店编号
	TerminalID   string `json:"terminal_id"`            // 可选 商户的终端编号
}

//APIName APIName
func (t TradeRefund) APIName() string {
	return "alipay.trade.refund"
}

//Params Params
func (t TradeRefund) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeRefund) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeRefund) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeRefundResponse TradeRefundResponse
type TradeRefundResponse struct {
	AliPayTradeRefund struct {
		Code                 string             `json:"code"`
		Msg                  string             `json:"msg"`
		SubCode              string             `json:"sub_code"`
		SubMsg               string             `json:"sub_msg"`
		TradeNo              string             `json:"trade_no"`                          // 支付宝交易号
		OutTradeNo           string             `json:"out_trade_no"`                      // 商户订单号
		BuyerLogonID         string             `json:"buyer_logon_id"`                    // 用户的登录id
		BuyerUserID          string             `json:"buyer_user_id"`                     // 买家在支付宝的用户id
		FundChange           string             `json:"fund_change"`                       // 本次退款是否发生了资金变化
		RefundFee            string             `json:"refund_fee"`                        // 退款总金额
		GmtRefundPay         string             `json:"gmt_refund_pay"`                    // 退款支付时间
		StoreName            string             `json:"store_name"`                        // 交易在支付时候的门店名称
		RefundDetailItemList []RefundDetailItem `json:"refund_detail_item_list,omitempty"` // 退款使用的资金渠道
	} `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"`
}

//IsSuccess IsSuccess
func (t *TradeRefundResponse) IsSuccess() bool {
	if t.AliPayTradeRefund.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//RefundDetailItem RefundDetailItem
type RefundDetailItem struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
}

//=======================================================================================================

// FastpayTradeRefundQuery 统一收单交易退款查询
// https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
type FastpayTradeRefundQuery struct {
	AppAuthToken string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	OutRequestNo string `json:"out_request_no"`         // 必须 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
}

//APIName APIName
func (t FastpayTradeRefundQuery) APIName() string {
	return "alipay.trade.fastpay.refund.query"
}

//Params Params
func (t FastpayTradeRefundQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t FastpayTradeRefundQuery) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t FastpayTradeRefundQuery) ExtJSONParamValue() string {
	return marshal(t)
}

//FastpayTradeRefundQueryResponse FastpayTradeRefundQueryResponse
type FastpayTradeRefundQueryResponse struct {
	AliPayTradeFastpayRefundQueryResponse struct {
		Code         string `json:"code"`
		Msg          string `json:"msg"`
		SubCode      string `json:"sub_code"`
		SubMsg       string `json:"sub_msg"`
		OutRequestNo string `json:"out_request_no"` // 本笔退款对应的退款请求号
		OutTradeNo   string `json:"out_trade_no"`   // 创建交易传入的商户订单号
		RefundReason string `json:"refund_reason"`  // 发起退款时，传入的退款原因
		TotalAmount  string `json:"total_amount"`   // 发该笔退款所对应的交易的订单金额
		RefundAmount string `json:"refund_amount"`  // 本次退款请求，对应的退款金额
		TradeNo      string `json:"trade_no"`       // 支付宝交易号
	} `json:"alipay_trade_fastpay_refund_query_response"`
	Sign string `json:"sign"`
}

// IsSuccess IsSuccess
func (t *FastpayTradeRefundQueryResponse) IsSuccess() bool {
	if t.AliPayTradeFastpayRefundQueryResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//=======================================================================================================

//TradeOrderSettle 统一收单交易结算接口
// https://docs.open.alipay.com/api_1/alipay.trade.order.settle
type TradeOrderSettle struct {
	AppAuthToken      string             `json:"-"`                  // 可选
	OutRequestNo      string             `json:"out_request_no"`     // 必须 结算请求流水号 开发者自行生成并保证唯一性
	TradeNo           string             `json:"trade_no"`           // 必须 支付宝订单号
	RoyaltyParameters []RoyaltyParameter `json:"royalty_parameters"` // 必须 分账明细信息
	OperatorID        string             `json:"operator_id"`        //可选 操作员id
}

//APIName APIName
func (t TradeOrderSettle) APIName() string {
	return "alipay.trade.order.settle"
}

//Params Params
func (t TradeOrderSettle) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

// ExtJSONParamName ExtJSONParamName
func (t TradeOrderSettle) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeOrderSettle) ExtJSONParamValue() string {
	return marshal(t)
}

//RoyaltyParameter RoyaltyParameter
type RoyaltyParameter struct {
	TransOut         string  `json:"trans_out"`         // 可选 分账支出方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransIn          string  `json:"trans_in"`          // 可选 分账收入方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	Amount           float64 `json:"amount"`            // 可选 分账的金额，单位为元
	AmountPercentage float64 `json:"amount_percentage"` // 可选 分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	Desc             string  `json:"desc"`              // 可选 分账描述
}

//TradeOrderSettleResponse TradeOrderSettleResponse
type TradeOrderSettleResponse struct {
	Body struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
		TradeNo string `json:"trade_no"`
	} `json:"alipay_trade_order_settle_response"`
	Sign string `json:"sign"`
}

//=======================================================================================================

// TradeCreate 统一收单交易创建接口
// https://docs.open.alipay.com/api_1/alipay.trade.create/
type TradeCreate struct {
	TradePay
	AppAuthToken       string            `json:"-"`                   // 可选
	DiscountableAmount string            `json:"discountable_amount"` // 可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位
	BuyerID            string            `json:"buyer_id"`
	GoodsDetail        []GoodsDetailItem `json:"goods_detail,omitempty"`
	OperatorID         string            `json:"operator_id"`
	TerminalID         string            `json:"terminal_id"`
}

//APIName APIName
func (t TradeCreate) APIName() string {
	return "alipay.trade.create"
}

//Params Params
func (t TradeCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeCreate) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeCreate) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeCreateResponse TradeCreateResponse
type TradeCreateResponse struct {
	AliPayTradeCreateResponse struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		SubCode    string `json:"sub_code"`
		SubMsg     string `json:"sub_msg"`
		TradeNo    string `json:"trade_no"` // 支付宝交易号
		OutTradeNo string `json:"out_trade_no"`
	} `json:"alipay_trade_create_response"`
	Sign string `json:"sign"`
}

//ExtendParamsItem ExtendParamsItem
type ExtendParamsItem struct {
	SysServiceProviderID string `json:"sys_service_provider_id"`
	HbFqNum              string `json:"hb_fq_num"`
	HbFqSellerPercent    string `json:"hb_fq_seller_percent"`
	TimeoutExpress       string `json:"timeout_express"`
}

//RoyaltyInfo RoyaltyInfo
type RoyaltyInfo struct {
	RoyaltyType       string                  `json:"royalty_type"`
	RoyaltyDetailInfo []RoyaltyDetailInfoItem `json:"royalty_detail_infos,omitempty"`
}

//RoyaltyDetailInfoItem RoyaltyDetailInfoItem
type RoyaltyDetailInfoItem struct {
	SerialNo         string `json:"serial_no"`
	TransInType      string `json:"trans_in_type"`
	BatchNo          string `json:"batch_no"`
	OutRelationID    string `json:"out_relation_id"`
	TransOutType     string `json:"trans_out_type"`
	TransOut         string `json:"trans_out"`
	TransIn          string `json:"trans_in"`
	Amount           string `json:"amount"`
	Desc             string `json:"desc"`
	AmountPercentage string `json:"amount_percentage"`
	AliPayStoreID    string `json:"alipay_store_id"`
}

//SubMerchantItem SubMerchantItem
type SubMerchantItem struct {
	MerchantID string `json:"merchant_id"`
}

//GoodsDetailItem GoodsDetailItem
type GoodsDetailItem struct {
	GoodsID       string `json:"goods_id"`
	AliPayGoodsID string `json:"alipay_goods_id"`
	GoodsName     string `json:"goods_name"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	GoodsCategory string `json:"goods_category"`
	Body          string `json:"body"`
	ShowURL       string `json:"show_url"`
}

//=======================================================================================================

// AliPayTradePay 统一收单交易支付接口
// https://docs.open.alipay.com/api_1/alipay.trade.pay/
type AliPayTradePay struct {
	TradePay
	AppAuthToken string `json:"-"` // 可选

	Scene    string `json:"scene"`     // 必须 支付场景 条码支付，取值：bar_code 声波支付，取值：wave_code, bar_code, wave_code
	AuthCode string `json:"auth_code"` // 必须 支付授权码

	BuyerID            string            `json:"buyer_id"` // 可选 家的支付宝用户id，如果为空，会从传入了码值信息中获取买家ID
	TransCurrency      string            `json:"trans_currency,omitempty"`
	SettleCurrency     string            `json:"settle_currency,omitempty"`
	DiscountableAmount string            `json:"discountable_amount,omitempty"` // 可选 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]。 如果该值未传入，但传入了【订单总金额】和【不可打折金额】，则该值默认为【订单总金额】-【不可打折金额】
	GoodsDetail        []GoodsDetailItem `json:"goods_detail,omitempty"`        // 可选 订单包含的商品列表信息，Json格式，其它说明详见商品明细说明
	OperatorID         string            `json:"operator_id,omitempty"`         // 可选 商户操作员编号
	TerminalID         string            `json:"terminal_id,omitempty"`         // 可选 商户机具终端编号
	AuthConfirmMode    string            `json:"auth_confirm_mode,omitempty"`
	TerminalParams     string            `json:"terminal_params,omitempty"`
}

//APIName APIName
func (t AliPayTradePay) APIName() string {
	return "alipay.trade.pay"
}

//Params Params
func (t AliPayTradePay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t AliPayTradePay) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t AliPayTradePay) ExtJSONParamValue() string {
	return marshal(t)
}

//AliPayTradePayResponse AliPayTradePayResponse
type AliPayTradePayResponse struct {
	AliPayTradePay struct {
		Code                string          `json:"code"`
		Msg                 string          `json:"msg"`
		SubCode             string          `json:"sub_code"`
		SubMsg              string          `json:"sub_msg"`
		BuyerLogonID        string          `json:"buyer_logon_id"`           // 买家支付宝账号
		BuyerPayAmount      string          `json:"buyer_pay_amount"`         // 买家实付金额，单位为元，两位小数。
		BuyerUserID         string          `json:"buyer_user_id"`            // 买家在支付宝的用户id
		CardBalance         string          `json:"card_balance"`             // 支付宝卡余额
		DiscountGoodsDetail string          `json:"discount_goods_detail"`    // 本次交易支付所使用的单品券优惠的商品优惠信息
		FundBillList        []FundBill      `json:"fund_bill_list,omitempty"` // 交易支付使用的资金渠道
		GmtPayment          string          `json:"gmt_payment"`
		InvoiceAmount       string          `json:"invoice_amount"`                // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		OutTradeNo          string          `json:"out_trade_no"`                  // 创建交易传入的商户订单号
		TradeNo             string          `json:"trade_no"`                      // 支付宝交易号
		PointAmount         string          `json:"point_amount"`                  // 积分支付的金额，单位为元，两位小数。
		ReceiptAmount       string          `json:"receipt_amount"`                // 实收金额，单位为元，两位小数
		StoreName           string          `json:"store_name"`                    // 发生支付交易的商户门店名称
		TotalAmount         string          `json:"total_amount"`                  // 发该笔退款所对应的交易的订单金额
		VoucherDetailList   []VoucherDetail `json:"voucher_detail_list,omitempty"` // 本交易支付时使用的所有优惠券信息
	} `json:"alipay_trade_pay_response"`
	Sign string `json:"sign"`
}

//IsSuccess IsSuccess
func (t *AliPayTradePayResponse) IsSuccess() bool {
	if t.AliPayTradePay.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//=======================================================================================================

// TradeAppPay app支付接口2.0
// https://docs.open.alipay.com/api_1/alipay.trade.app.pay/
type TradeAppPay struct {
	TradePay
	TimeExpire string `json:"time_expire,omitempty"` // 绝对超时时间，格式为yyyy-MM-dd HH:mm。
}

//APIName APIName
func (t TradeAppPay) APIName() string {
	return "alipay.trade.app.pay"
}

//Params Params
func (t TradeAppPay) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeAppPay) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeAppPay) ExtJSONParamValue() string {
	return marshal(t)
}

//=======================================================================================================

//TradePreCreate  (统一收单线下交易预创建
// https://docs.open.alipay.com/api_1/alipay.trade.precreate/
type TradePreCreate struct {
	TradePay
	AppAuthToken       string            `json:"-"`                      // 可选
	DiscountableAmount string            `json:"discountable_amount"`    // 可选 可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000] 如果该值未传入，但传入了【订单总金额】，【不可打折金额】则该值默认为【订单总金额】-【不可打折金额】
	GoodsDetail        []GoodsDetailItem `json:"goods_detail,omitempty"` // 可选 订单包含的商品列表信息.Json格式. 其它说明详见：“商品明细说明”
	OperatorID         string            `json:"operator_id"`            // 可选 商户操作员编号
	TerminalID         string            `json:"terminal_id"`            // 可选 商户机具终端编号
}

//APIName APIName
func (t TradePreCreate) APIName() string {
	return "alipay.trade.precreate"
}

//Params Params
func (t TradePreCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradePreCreate) ExtJSONParamName() string {
	return "biz_content"
}

// ExtJSONParamValue ExtJSONParamValue
func (t TradePreCreate) ExtJSONParamValue() string {
	return marshal(t)
}

//TradePreCreateResponse TradePreCreateResponse
type TradePreCreateResponse struct {
	AliPayPreCreateResponse struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		SubCode    string `json:"sub_code"`
		SubMsg     string `json:"sub_msg"`
		OutTradeNo string `json:"out_trade_no"` // 创建交易传入的商户订单号
		QRCode     string `json:"qr_code"`      // 当前预下单请求生成的二维码码串，可以用二维码生成工具根据该码串值生成对应的二维码
	} `json:"alipay_trade_precreate_response"`
	Sign string `json:"sign"`
}

//IsSuccess IsSuccess
func (t *TradePreCreateResponse) IsSuccess() bool {
	if t.AliPayPreCreateResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//=======================================================================================================

// TradeCancel 统一收单交易撤销接口
// https://docs.open.alipay.com/api_1/alipay.trade.cancel/
type TradeCancel struct {
	AppAuthToken string `json:"-"` // 可选
	NotifyURL    string `json:"-"` // 可选

	OutTradeNo string `json:"out_trade_no"` // 原支付请求的商户订单号,和支付宝交易号不能同时为空
	TradeNo    string `json:"trade_no"`     // 支付宝交易号，和商户订单号不能同时为空
}

//APIName APIName
func (t TradeCancel) APIName() string {
	return "alipay.trade.cancel"
}

//Params Params
func (t TradeCancel) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	m["notify_url"] = t.NotifyURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeCancel) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeCancel) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeCancelResponse TradeCancelResponse
type TradeCancelResponse struct {
	AliPayTradeCancelResponse struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		SubCode    string `json:"sub_code"`
		SubMsg     string `json:"sub_msg"`
		TradeNo    string `json:"trade_no"`     // 支付宝交易号
		OutTradeNo string `json:"out_trade_no"` // 创建交易传入的商户订单号
		RetryFlag  string `json:"retry_flag"`   // 是否需要重试
		Action     string `json:"action"`       // 本次撤销触发的交易动作 close：关闭交易，无退款 refund：产生了退款
	} `json:"alipay_trade_cancel_response"`
	Sign string `json:"sign"`
}

//IsSuccess IsSuccess
func (t *TradeCancelResponse) IsSuccess() bool {
	if t.AliPayTradeCancelResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

//=======================================================================================================

// TradeOrderInfoSync 支付宝订单信息同步接口
// https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync/
type TradeOrderInfoSync struct {
	AppAuthToken string `json:"-"`              // 可选
	OutRequestNo string `json:"out_request_no"` // 必选 标识一笔交易多次请求，同一笔交易多次信息同步时需要保证唯一
	BizType      string `json:"biz_type"`       // 必选 交易信息同步对应的业务类型，具体值与支付宝约定；信用授权场景下传CREDIT_AUTH
	TradeNo      string `json:"trade_no"`       // 可选 支付宝交易号，和商户订单号不能同时为空
	OrderBizInfo string `json:"order_biz_info"` // 可选 商户传入同步信息，具体值要和支付宝约定；用于芝麻信用租车、单次授权等信息同步场景，格式为json格式
}

//APIName APIName
func (t TradeOrderInfoSync) APIName() string {
	return "alipay.trade.orderinfo.sync"
}

//Params Params
func (t TradeOrderInfoSync) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t TradeOrderInfoSync) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t TradeOrderInfoSync) ExtJSONParamValue() string {
	return marshal(t)
}

//TradeOrderInfoSyncResponse TradeOrderInfoSyncResponse
type TradeOrderInfoSyncResponse struct {
	Body struct {
		Code        string `json:"code"`
		Msg         string `json:"msg"`
		SubCode     string `json:"sub_code"`
		SubMsg      string `json:"sub_msg"`
		TradeNo     string `json:"trade_no"`
		OutTradeNo  string `json:"out_trade_no"`
		BuyerUserID string `json:"buyer_user_id"`
	} `json:"alipay_trade_orderinfo_sync_response"`
	Sign string `json:"sign"`
}
