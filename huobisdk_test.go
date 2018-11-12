package robot

import (
	"log"
	"testing"
	"time"
)

func getHuobiSdk() *HuobiSdk {
	h := &HuobiSdk{
		Host:        "api.huobi.pro",
		AccessKeyID: "your access id",
		SecretKey:   "your secret id",
	}
	return h
}

func TestGetAccounts(t *testing.T) {
	res, err := getHuobiSdk().Accounts()
	if err != nil {
		log.Fatalf("failed:%v", err)
	} else {
		log.Printf("get huobi result:%v", res)
	}

}

func TestGetMarketDepth(t *testing.T) {
	cc, err := getHuobiSdk().MarketDepth("your symbol")
	if err != nil {
		log.Fatalf("fail get depth:%v", err)
	} else {
		log.Printf("get depth: %v", cc.Tick.Bids[0])
	}
}

func TestOrderPlace(t *testing.T) {
	cc, err := getHuobiSdk().OrderPlace("your symbol", "price", "amount", "buy-limit")
	if err != nil {
		log.Fatalf("fail buy:%v", err)
	} else {
		log.Printf("下单成功，订单号为:%v", cc.Data)
	}
}

func TestOrderDetail(t *testing.T) {
	cc, err := getHuobiSdk().OrderDetail("15298248443")
	if err != nil {
		log.Fatalf("查看订单详细信息失败:%v", err)
	} else {
		log.Printf("查看信息成功:%v", cc.Data.State)
	}
}

func TestOrderCancel(t *testing.T) {
	cc, err := getHuobiSdk().OrderCancel("15298248443")
	if err != nil {
		log.Fatalf("取消订单失败:%v", err)
	} else {
		log.Printf("取消订单成功:%v", cc.Status)
	}

}

func ppp() interface{} {
	return OrderResponse{Status: "ok", Data: "some"}
}

func TestBasic(t *testing.T) {
	hbsk := getHuobiSdk()
	res, err := hbsk.CurOrders("pnteth", "buy-limit", "submitted%2Cfilled")
	if err != nil {
		log.Fatalf("err is : %v", err)
	} else {
		log.Printf("current orders are:%v", len(res.Data))
	}
}

func TestAsync(t *testing.T) {
	dp := asyncGet()
	for _, v := range dp {
		log.Printf("结果:%v", v.dr.Tick.Ts)
	}
}

type depthErr struct {
	dr  *DepthResponse
	err error
}

func asyncGet() []*depthErr {
	ch := make(chan *depthErr)
	responses := []*depthErr{}
	for i := 0; i < 10; i++ {
		go func() {
			d, err := getHuobiSdk().MarketDepth("pnteth")
			ch <- &depthErr{d, err}
		}()
	}
loop:
	for {
		select {
		case r := <-ch:
			responses = append(responses, r)
			if len(responses) == 10 {
				break loop
			}
		case <-time.After(5 * time.Second):
			log.Printf("请求")
		}
	}
	return responses
}

func TestOrdersCancel(t *testing.T) {
	h := getHuobiSdk()

	res, err := h.OrdersCancel([]string{"1"})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("res is:%v", res)
	}

}
