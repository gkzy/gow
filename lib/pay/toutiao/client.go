package toutiao

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/clbanning/mxj"
	"github.com/gkzy/gow/lib/logy"
	"io/ioutil"
	"net/http"
	"crypto/md5"
	"sort"
	"strings"
)

const (
	bodyType       = "application/json"
	createOrderUrl = "https://developer.toutiao.com/api/apps/ecpay/v1/create_order"
	queryOrderUrl  = "https://developer.toutiao.com/api/apps/ecpay/v1/query_order"
)

type Client struct {
	AppId     string //传入的appID
	NotifyURL string //异步通知地址
	OrderTime int    //订单有效分钟数
	SALT      string //向小程序平台发送请求时的密钥
	Token     string //小程序平台向开发者服务端发送请求时的密钥
	Extra     string //开发者自定义字段，回调原样回传
}

//NewClient 一个新的配置信息
func NewClient(appId, SALT, token string, notifyUrl string, orderTime int) *Client {
	return &Client{
		AppId:     appId,
		NotifyURL: notifyUrl,
		OrderTime: orderTime,
		SALT:      SALT,
		Token:     token,
	}
}

//预下单
func (c *Client) CreateOrder(body, outTradeNo string, totalFee int64) (rslt *OrderInfo, err error) {
	params := make(Params)
	params.SetString("out_order_no", outTradeNo)
	params.SetInt64("total_amount", totalFee)
	params.SetString("subject", body)
	params.SetString("body", body)
	params.SetInt64("valid_time", int64(c.OrderTime*60)) //订单过期时间(秒); 最小 15 分钟，最大两天
	params.SetString("cp_extra", c.Extra)                 //开发者自定义字段，回调原样回传
	params.SetString("notify_url", c.NotifyURL)
	params.SetInt64("disable_msg", int64(0)) //是否屏蔽担保支付的推送消息，1-屏蔽 0-非屏蔽，

	resp, err := c.post(createOrderUrl, params)
	if err != nil {
		logy.Errorf("下单出错:%v", err)
		return
	}
	data := new(CreateOrderResp)
	if err = json.Unmarshal([]byte(resp), &data); err != nil {
		logy.Errorf("解析预下单响应参数出错:%v", err)
		return
	}
	//状态码 0-业务处理成功
	if data.ErrNo == 0 {
		rslt = data.Data
	} else {
		logy.Errorf("预下单业务处理失败,错误码：%v,错误信息:%v", data.ErrNo, data.ErrTips)
		return
	}
	return
}

//订单查询
func (c *Client) QueryOrder(outTradeNo string) (rslt *QueryOrderRespData, err error) {
	params := make(Params)
	params.SetString("out_order_no", outTradeNo)
	resp, err := c.post(queryOrderUrl, params)
	if err != nil {
		logy.Errorf("查订单出错:%v", err)
		return
	}
	//fmt.Println("resp:",resp)
	ret := new(QueryOrderResp)
	err = json.Unmarshal([]byte(resp), &ret)
	if err != nil {
		logy.Errorf("解析查询订单响应参数出错:%v", err)
		return
	}
	rslt = ret.PaymentInfo
	return
}

//回调
func (c *Client) Notify(req *http.Request) (msgData *NotifyMsgData, ret *NotifyReturn, err error) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	// 写回 body 内容
	req.Body = ioutil.NopCloser(bytes.NewReader(content))
	params := StringToMap(content)
	var returntype string
	if params.ContainsKey("type") {
		returntype = params.GetString("type")
	} else {
		return nil, nil, errors.New("没有回调类型标记")
	}
	//回调类型标记，支付成功回调为"payment"
	if returntype == "payment" {
		//验证签名
		if c.ValidSign(params) {
			msg := params.GetString("msg")
			msgData = new(NotifyMsgData)
			json.Unmarshal([]byte(msg), &msgData)
			ret = new(NotifyReturn)
			ret.ErrNo = 0
			ret.ErrTip = "success"
			return
		}
		return nil, nil, errors.New("签名不正确")
	} else {
		return nil, nil, errors.New("回调类型不为支付成功")
	}
}

// ValidSign 验证签名
func (c *Client) ValidSign(params Params) bool {
	if !params.ContainsKey("msg_signature") {
		return false
	}
	return params.GetString("msg_signature") == c.respsign(params)
}

func (c *Client) post(url string, params Params) (string, error) {
	h := &http.Client{}
	p := c.fullRequestParams(params)
	response, err := h.Post(url, bodyType, strings.NewReader(MapToString(p)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *Client) fullRequestParams(params Params) Params {
	params["app_id"] = c.AppId
	signstr := c.reqsign(params)
	//fmt.Println("得到的签名字符串：", signstr)
	params["sign"] = signstr
	return params
}

//请求签名
func (c *Client) reqsign(params Params) string {
	// 创建切片
	var values = make([]string, 0, len(params))
	values = append(values, c.SALT)
	// 遍历签名参数
	for k, v := range params {
		if k == "sign" || k == "app_id" || k == "thirdparty_id" { //sign, app_id , thirdparty_id 字段用于标识身份字段，不参与签名
			continue
		}
		values = append(values, fmt.Sprintf("%v", v))
	}
	//fmt.Println("签名前value:", values)
	sort.Strings(values)
	h := md5.New()
	h.Write([]byte(strings.Join(values, "&")))
	var toSignStr string
	return fmt.Sprintf("%x", h.Sum([]byte(toSignStr)))
}

//回调签名
func (c *Client) respsign(params Params) string {
	// 创建切片
	var values = make([]string, 0, len(params))
	values = append(values, c.Token)
	// 遍历签名参数
	for k, v := range params {
		if k == "msg_signature" || k == "type" { //验证时注意不包含 msg_signature 签名本身，不包含空字段与 type 常量字段
			continue
		}
		values = append(values, fmt.Sprintf("%v", v))
	}
	//fmt.Println("签名前value:", values)
	sort.Strings(values)
	h := sha1.New()
	h.Write([]byte(strings.Join(values, "")))
	bs := h.Sum(nil)
	_signature := fmt.Sprintf("%x", bs)
	return _signature
}

func MapToString(params Params) string {
	data, _ := json.Marshal(params)
	return string(data)
}

func StringToMap(content []byte) Params {
	param := make(Params)
	json.Unmarshal(content, &param)
	return param
}
