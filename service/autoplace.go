package service

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ryzencool/huobisdk-go"
	"github.com/shopspring/decimal"
)

// StepList pl
type StepList []robot.Step

func (s StepList) Len() int {
	return len(s)
}

func (s StepList) Less(i, j int) bool {
	return s[i].FlagEnd*s[i].FlagStart < s[j].FlagEnd*s[j].FlagStart
}

func (s StepList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// DecimalList decimal list implement sort interface
type DecimalList []decimal.Decimal

func (d DecimalList) Len() int {
	return len(d)
}

func (d DecimalList) Less(i, j int) bool {
	return d[i].Cmp(d[j]) < 0
}

func (d DecimalList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type order struct {
	id    int
	price decimal.Decimal
}

// AutoPlace is auto place order
func AutoPlace(hsdk *robot.HuobiSdk, pc robot.Place) {
	symbol := pc.Symbol
	rand.Seed(time.Now().UnixNano())
	md, err := hsdk.MarketDepth(symbol)
	if err != nil {
		return
	}
	buyMd := md.Tick.Bids
	sellMd := md.Tick.Asks
	if len(buyMd) == 0 || len(sellMd) == 0 {
		return
	}
	var ps []float64
	for _, v := range buyMd {
		ps = append(ps, v[0])
	}
	prec, err := getPrec(ps)
	if err != nil {
		return
	}
	buy0MdDP := decimal.NewFromFloat(buyMd[0][0])
	sell0MdDP := decimal.NewFromFloat(sellMd[0][0])

	var buyDpt20 []decimal.Decimal
	for _, b := range buyMd[:19] {
		v := decimal.NewFromFloat(b[0])
		buyDpt20 = append(buyDpt20, v)
	}

	var sellDpt20 []decimal.Decimal

	for _, s := range sellMd[:19] {
		v := decimal.NewFromFloat(s[0])
		sellDpt20 = append(sellDpt20, v)
	}

	buy1to5, err := continuousNums(buy0MdDP, prec, 5, "-")
	sell1to5, err := continuousNums(sell0MdDP, prec, 5, "+")
	b6, err := decimal.NewFromString(buy1to5[4])
	if err != nil {
		return
	}
	s6, err := decimal.NewFromString(sell1to5[4])
	if err != nil {
		return
	}
	buy6to20, err := randNums(b6, prec, 15, 2, 5, "-")
	if err != nil {
		return
	}
	sell6to20, err := randNums(s6, prec, 15, 2, 5, "+")
	if err != nil {
		return
	}
	bmerge := mergeSlice(buy1to5, buy6to20)
	smerge := mergeSlice(sell1to5, sell6to20)

	bset := make(map[string]struct{})
	sset := make(map[string]struct{})

	var bmergeDecimal []decimal.Decimal
	var smergeDecimal []decimal.Decimal
	for _, bs := range bmerge {
		tmp, err := decimal.NewFromString(bs)
		if err != nil {
			return
		}
		bmergeDecimal = append(bmergeDecimal, tmp)
	}

	for _, ss := range smerge {
		sd, err := decimal.NewFromString(ss)
		if err != nil {
			return
		}
		smergeDecimal = append(smergeDecimal, sd)
	}

	for _, v := range buyDpt20 {
		r := v.String()
		bset[r] = struct{}{}
	}

	for _, v := range bmergeDecimal {
		r := v.String()
		bset[r] = struct{}{}
	}

	for _, v := range sellDpt20 {
		r := v.String()
		sset[r] = struct{}{}
	}

	for _, v := range smergeDecimal {
		r := v.String()
		sset[r] = struct{}{}
	}
	var bcv []string
	bc := sortSet(bset, -1)
	for _, v := range bc {
		if v < "1" {
			bcv = append(bcv, v)
		}
	}
	var scv []string
	sc := sortSet(sset, 1)
	for _, v := range sc {
		if v < "1" {
			scv = append(scv, v)
		}
	}

	bc20 := bcv[:20]
	sc20 := scv[:20]
	bc20map := make(map[string]string)
	sc20map := make(map[string]string)
	for _, ps := range pc.Step {
		for i := ps.FlagStart - 1; i < ps.FlagEnd; i++ {
			bc20map[bc20[i]] = findRand(ps.AmountMin, ps.AmountMax)
			sc20map[sc20[i]] = findRand(ps.AmountMin, ps.AmountMax)
		}
	}

	cbs, err := hsdk.CurOrders(symbol, "buy-limit", "submitted%2Cpartial-filled")
	if err != nil || cbs.Status != "ok" {
		log.Printf("获取买委托失败:%v, %v", err, cbs)
		return
	}
	for _, v := range cbs.Data {
		r, err := decimal.NewFromString(v.Price)

		if err != nil {
			return
		}
		rs := r.String()
		for i := 0; i < len(bc20); i++ {
			if rs == bc20[i] {
				bc20 = append(bc20[:i], bc20[i+1:]...)
				i--
			}
		}
	}

	css, err := hsdk.CurOrders(symbol, "sell-limit", "submitted%2Cpartial-filled")
	if err != nil || css.Status != "ok" {
		return
	}
	// log.Printf("卖委托单列表:%v", css.Data)
	for _, v := range css.Data {
		r, err := decimal.NewFromString(v.Price)
		if err != nil {
			return
		}
		rs := r.String()
		for i := 0; i < len(sc20); i++ {
			if rs == sc20[i] {
				sc20 = append(sc20[:i], sc20[i+1:]...)
				i--
			}
		}
	}

	log.Printf("<买单>:%v", bc20)
	log.Printf("<卖单>:%v", sc20)

	for _, v := range bc20 {
		bm := bc20map[v]
		log.Printf("[铺单-%v] 下买单 价格:%v, 数量:%v", symbol, v, bm)
		_, err := hsdk.OrderPlace(symbol, v, bm, "buy-limit")
		if err != nil {
			return
		}
	}

	for _, v := range sc20 {
		sm := sc20map[v]
		_, err := hsdk.OrderPlace(symbol, v, sm, "sell-limit")
		if err != nil {
			return
		}
	}

}

// CancelOver30 cancel order over 30
func CancelOver30(hbsk *robot.HuobiSdk, symbol string) {

	md, err := hbsk.MarketDepth(symbol)
	var buyList []order
	var sellList []order
	if err != nil {
		return
	}
	buyMd := md.Tick.Bids
	sellMd := md.Tick.Asks
	if len(buyMd) < 30 {
		return
	}

	appendList := func(org *robot.CurOrdersResponse, list []order) ([]order, error) {
		if len(org.Data) > 0 {
			for _, v := range org.Data {
				r, err := decimal.NewFromString(v.Price)
				if err != nil {
					return nil, err
				}
				res := order{id: v.ID, price: r}
				list = append(list, res)
			}
		}
		return list, nil
	}
	bs, err := hbsk.CurOrders(symbol, "buy-limit", "submitted%2Cpartial-filled")
	if err != nil {
		return
	}
	buyList, _ = appendList(bs, buyList)
	sp, err := hbsk.CurOrders(symbol, "sell-limit", "submitted%2Cpartial-filled")
	if err != nil {
		return
	}
	sellList, _ = appendList(sp, sellList)
	var buyCancelList []string
	var sellCancelList []string
	if len(buyMd) > 30 {
		buy30Pr := decimal.NewFromFloat(buyMd[29][0])

		for _, v := range buyList {
			if v.price.Cmp(buy30Pr) < 0 {
				buyCancelList = append(buyCancelList, strconv.Itoa(v.id))
			}
		}
		go hbsk.OrdersCancel(buyCancelList)
	}
	if len(sellMd) > 30 {
		sell30Pr := decimal.NewFromFloat(sellMd[29][0])
		for _, v := range sellList {
			if v.price.Cmp(sell30Pr) > 0 {
				sellCancelList = append(sellCancelList, strconv.Itoa(v.id))
			}
		}
		go hbsk.OrdersCancel(sellCancelList)
	}
}

// get majarioty length of float slice
func getPrec(fa []float64) (int, error) {
	s := make(map[int]int)
	for _, f := range fa {
		df := decimal.NewFromFloat(f)
		dl := len(df.String())
		if _, ok := s[dl]; ok {
			s[dl]++
		} else {
			s[dl] = 1
		}
	}
	var t []int
	for _, v := range s {
		t = append(t, v)
	}
	sort.Ints(t)
	max := t[len(t)-1]
	for k, v := range s {
		if v == max {
			return k, nil
		}
	}
	return -1, nil
}

// appropriate to decimal
func calDecimal(org decimal.Decimal, l int, n int, orderType string) (string, error) {
	orgStr := org.String()
	lorg := len(orgStr)
	for i := 0; i < l-lorg; i++ {
		orgStr += "0"
	}
	orgStr = strings.TrimLeftFunc(orgStr, func(r rune) bool {
		return r == '0' || r == '.'
	})
	d, err := decimal.NewFromString(orgStr)
	if err != nil {
		return "", err
	}
	var da string
	var lr int
	if orderType == "+" {
		da = d.Add(decimal.New(int64(n), 0)).String()
		lr = l - len(orgStr) - 2 - (len(da) - len(orgStr))
	} else if orderType == "-" {
		da = d.Sub(decimal.New(int64(n), 0)).String()
		lr = l - len(orgStr) - 2 + (len(orgStr) - len(da))
	} else {
		return "", fmt.Errorf("类型参数错误")
	}

	for i := 0; i < lr; i++ {
		da = "0" + da
	}
	return "0." + da, nil

}

func continuousNums(org decimal.Decimal, precision int, amount int, tp string) ([]string, error) {
	list := make([]string, amount)
	for i := 0; i < amount; i++ {
		rs, err := calDecimal(org, precision, i, tp)
		if err != nil {
			return nil, err
		}
		list[i] = rs
	}
	return list, nil
}

func randNums(org decimal.Decimal, precision int, amount int, min int, max int, tp string) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	list := make([]string, amount)
	start := org
	for i := 0; i < amount; i++ {
		rn := min + rand.Intn(max-min)
		rs, err := calDecimal(start, precision, rn, tp)
		if err != nil {
			return nil, err
		}
		start, err = decimal.NewFromString(rs)
		list[i] = start.String()
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func mergeSlice(sls ...[]string) []string {
	list := make([]string, 0)
	for _, vv := range sls {
		for _, v := range vv {
			list = append(list, v)
		}
	}
	return list
}

type orderPlaceReponseErr struct {
	op  *robot.OrderResponse
	err error
}

func asyncPlaceOrder(hbsk *robot.HuobiSdk, symbol string, price string, amount string, tp string) []*orderPlaceReponseErr {
	ch := make(chan *orderPlaceReponseErr)
	responses := []*orderPlaceReponseErr{}
	for i := 0; i < 10; i++ {
		go func() {
			d, err := hbsk.OrderPlace(symbol, price, amount, tp)
			ch <- &orderPlaceReponseErr{d, err}
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

func sortSet(org map[string]struct{}, tp int) []string {
	var sList []string
	for k := range org {
		sList = append(sList, k)
	}
	if tp >= 0 {
		sort.Strings(sList)

	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(sList)))
	}
	return sList
}
