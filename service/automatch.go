package service

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ryzencool/huobisdk-go"
	"github.com/shopspring/decimal"
)

type orderResponseErr struct {
	or  *robot.OrderResponse
	err error
}

var ch = make(chan orderResponseErr)
var ch0 = make(chan orderResponseErr)

// AutoMatch auto match
func AutoMatch(hsdk *robot.HuobiSdk, symbol string, limit int64, amountMin int, amountMax int) {

	depth, err := hsdk.MarketDepth(symbol)
	if err != nil {

		return
	}

	buy0 := depth.Tick.Bids[0]
	sell0 := depth.Tick.Asks[0]
	buy0Price := buy0[0]
	sell0Price := sell0[0]

	buy0PriceD := decimal.NewFromFloat(buy0Price)
	sell0PriceD := decimal.NewFromFloat(sell0Price)
	numNear := checkNumNear(buy0PriceD.String(), sell0PriceD.String())
	if numNear {

		return
	}
	fr := findRand(float64(amountMin), float64(amountMax))
	done := make(chan interface{})
	defer close(done)
	var chs orderResponseErr
	var ch0s orderResponseErr

	pr := findRand(buy0Price, sell0Price)
	placeOrder(done, hsdk, symbol, pr, fr, ch, ch0)
	chs = <-ch
	ch0s = <-ch0
	if chs.err != nil || chs.or.Status != "ok" {
		return
	}
	if err != nil || ch0s.or.Status != "ok" {
		return
	}

	go cancelOrder(hsdk, chs.or.Data, symbol)
	go cancelOrder(hsdk, ch0s.or.Data, symbol)
}

func checkNumNear(buy, sell string) bool {
	fn := func(r rune) bool {
		return r == '.' || r == '0'
	}
	buyW := strings.TrimLeftFunc(buy, fn)
	sellW := strings.TrimLeftFunc(sell, fn)
	bl := len(buyW)
	sl := len(sellW)
	if bl > sl {
		n := bl - sl
		for i := 0; i < n; i++ {
			sellW += "0"
		}
	} else if bl < sl {
		n := sl - bl
		for i := 0; i < n; i++ {
			buyW += "0"
		}
	}
	bb, _ := decimal.NewFromString(buyW)
	ss, _ := decimal.NewFromString(sellW)
	ee := decimal.New(1, 0)
	return ss.Sub(bb).Equal(ee)
}

func findRand(min, max float64) string {
	rand.Seed(time.Now().UnixNano())

	bp := decimal.NewFromFloat(min)
	sp := decimal.NewFromFloat(max)
	r := decimal.NewFromFloat(rand.Float64())
	res := (sp.Sub(bp)).Mul(r).Add(bp)
	bl := len(bp.String())
	sl := len(sp.String())
	var l int
	if bl >= sl {
		l = bl
	} else {
		l = sl
	}
	rs := res.String()[:l]
	rsd, _ := decimal.NewFromString(rs)
	if rsd.String() == bp.String() || rsd.String() == sp.String() {
		return findRand(min, max)
	}
	return rs

}

func cancelOrder(hsdk *robot.HuobiSdk, data string, symbol string) {
	odb, err := hsdk.OrderDetail(data)
	if err != nil {
		return
	}

	if odb.Data.State == "submitting" || odb.Data.State == "submitted" || odb.Data.State == "partial-filled" {
		_, err := hsdk.OrderCancel(strconv.Itoa(odb.Data.ID))
		if err != nil {
			return
		}
	}
}

func placeOrder(done <-chan interface{}, hsdk *robot.HuobiSdk, symbol string, price string, amount string, bch chan orderResponseErr, sch chan orderResponseErr) {
	go func() {
		orb, err := hsdk.OrderPlace(symbol, price, amount, "buy-limit")
		for {
			select {
			case <-done:
				return
			case bch <- orderResponseErr{orb, err}:
			}
		}
	}()
	go func() {
		ors, err := hsdk.OrderPlace(symbol, price, amount, "sell-limit")
		for {
			select {
			case <-done:
				return
			case sch <- orderResponseErr{ors, err}:
			}
		}
	}()
}
