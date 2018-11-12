package robot

import (
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// HuobiSdk is connect huobi api
type HuobiSdk struct {
	AccessKeyID string
	SecretKey   string
	Host        string
}

// Signature is huobi
func (h *HuobiSdk) signature(method string, webURL string, params map[string]string) string {
	params["AccessKeyId"] = h.AccessKeyID
	params["SignatureVersion"] = "2"
	params["SignatureMethod"] = "HmacSHA256"
	u := time.Now().UTC()
	f := strings.Join(strings.Split(strings.Split(u.String(), ".")[0], " "), "T")
	params["Timestamp"] = url.QueryEscape(f)
	var s []string
	for k := range params {
		k = url.QueryEscape(k)
		s = append(s, k)
	}
	sort.Strings(s)
	var s1 []string
	for _, v := range s {
		s1 = append(s1, v+"="+params[v])
	}
	p := strings.Join(s1, "&")
	origin := method + "\n" + h.Host + "\n" + webURL + "\n" + p
	params["Signature"] = url.QueryEscape(computeHmac256(origin, h.SecretKey))
	var s2 []string
	for k, v := range params {
		s2 = append(s2, k+"="+v)
	}
	return strings.Join(s2, "&")
}

// Accounts get user's account
func (h *HuobiSdk) Accounts() (*AccountsResponse, error) {
	var params = make(map[string]string)
	webURL := "/v1/account/accounts"
	h.signature("GET", webURL, params)
	hostURL := "https://" + h.Host + webURL
	res, err := Get(hostURL, params)
	if err != nil {
		return nil, err
	}
	ar := &AccountsResponse{}
	err = json.Unmarshal([]byte(res), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

// MarketDepth get market depth
func (h *HuobiSdk) MarketDepth(symbol string) (*DepthResponse, error) {
	var params = map[string]string{"symbol": symbol, "type": "step0"}
	webURL := "/market/depth"
	hostURL := "https://" + h.Host + webURL
	res, err := Get(hostURL, params)
	if err != nil {
		return nil, err
	}
	dr := &DepthResponse{}
	err = json.Unmarshal([]byte(res), dr)
	if err != nil {
		return nil, err
	}
	return dr, nil
}

// OrderPlace place order, source only support api
func (h *HuobiSdk) OrderPlace(symbol string, price string, amount string, orderType string) (*OrderResponse, error) {
	// log.Printf("<下单>: 交易对:%v, 价格:%v, 数量:%v, 类型:%v", symbol, price, amount, orderType)
	ar, err := h.Accounts()
	if err != nil {
		return nil, err
	}
	var id int
	if ar.Status == "ok" && len(ar.Data) > 0 {
		id = ar.Data[0].ID
	}
	accountID := strconv.Itoa(id)
	var params = map[string]string{
		"account-id": accountID,
		"amount":     amount,
		"price":      price,
		"source":     "api",
		"symbol":     symbol,
		"type":       orderType,
	}
	webURL := "/v1/order/orders/place"
	var query = make(map[string]string)
	q := h.signature("POST", webURL, query)
	hostURL := "https://" + h.Host + webURL + "?" + q
	tmp := make(map[string]interface{})
	for k, v := range params {
		tmp[k] = v
	}
	res, err := Post(hostURL, tmp)
	if err != nil {
		return nil, err
	}
	opr := &OrderResponse{}
	err = json.Unmarshal([]byte(res), opr)
	if err != nil {
		return nil, err
	}
	return opr, nil
}

// OrderDetail get some order's detail
func (h *HuobiSdk) OrderDetail(orderID string) (*OrderDetailResponse, error) {
	params := make(map[string]string)
	webURL := "/v1/order/orders/" + orderID
	h.signature("GET", webURL, params)
	hostURL := "https://" + h.Host + webURL
	res, err := Get(hostURL, params)
	if err != nil {
		return nil, err
	}
	odr := &OrderDetailResponse{}
	err = json.Unmarshal([]byte(res), odr)
	if err != nil {
		return nil, err
	}
	return odr, nil
}

// OrderCancel cancel some order
func (h *HuobiSdk) OrderCancel(orderID string) (*OrderResponse, error) {
	webURL := "/v1/order/orders/" + orderID + "/submitcancel"
	params := make(map[string]string)
	query := h.signature("POST", webURL, params)
	hostURL := "https://" + h.Host + webURL + "?" + query
	tmp := make(map[string]interface{})
	for k, v := range params {
		tmp[k] = v
	}
	res, err := Post(hostURL, tmp)
	if err != nil {
		// log.Fatal("取消订单失败")
		return nil, err
	}
	or := &OrderResponse{}
	err = json.Unmarshal([]byte(res), or)
	if err != nil {
		return nil, err
	}
	return or, nil
}

// OrdersCancel cancel orders
func (h *HuobiSdk) OrdersCancel(ids []string) (string, error) {
	webURL := "/v1/order/orders/batchcancel"
	paramsSign := make(map[string]string)

	query := h.signature("POST", webURL, paramsSign)
	hostURL := "https://" + h.Host + webURL + "?" + query
	params := make(map[string]interface{})
	params["order-ids"] = ids
	res, err := Post(hostURL, params)
	if err != nil {
		return "", err
	}

	return res, nil
}

// CurOrders get current orders
func (h *HuobiSdk) CurOrders(symbol string, tp string, states string) (*CurOrdersResponse, error) {
	webURL := "/v1/order/orders"
	params := make(map[string]string)
	params["symbol"] = symbol
	params["types"] = tp
	params["states"] = states
	params["size"] = "500"
	h.signature("GET", webURL, params)
	hostURL := "https://" + h.Host + webURL
	res, err := Get(hostURL, params)
	if err != nil {
		return nil, err
	}
	cor := &CurOrdersResponse{}
	err = json.Unmarshal([]byte(res), cor)
	if err != nil {
		return nil, err
	}
	return cor, nil
}
