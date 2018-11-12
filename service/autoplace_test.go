package service

import (
	"log"
	"sort"
	"strconv"
	"testing"

	"github.com/robfig/cron"
	"github.com/ryzencool/huobisdk-go"
	"github.com/shopspring/decimal"
)

func TestAutoPlace(t *testing.T) {
	hbsk := getHuobiSDK()
	pc := getPlaceConfig()
	go func() {
		for {
			AutoPlace(hbsk, pc)
		}
	}()

	for i := 0; i < 5; i++ {
		go func() {
			for {
				CancelOver30(hbsk, pc.Symbol)
			}
		}()
	}
	select {}
}

func TestCancelOver30(t *testing.T) {
	hbsk := getHuobiSDK()
	pc := getPlaceConfig()
	CancelOver30(hbsk, pc.Symbol)
	select {}
}

func TestCancel(t *testing.T) {
	h := getHuobiSDK()
	pc := getPlaceConfig()
	c := cron.New()
	c.AddFunc("*/3 * * * * ?", func() {
		cancel(h, pc)
	})
	c.Start()
	select {}
}

func cancel(h *robot.HuobiSdk, pc robot.Place) {

	bs, err := h.CurOrders(pc.Symbol, "buy-limit", "submitted%2Cpartial-filled")
	if err != nil {
		log.Printf("fail buy:%v", err)
	}
	var bl []string
	for _, v := range bs.Data {
		bl = append(bl, strconv.Itoa(v.ID))
	}
	log.Printf("fail buy:%v", bl)
	var sl []string
	ss, err := h.CurOrders(pc.Symbol, "sell-limit", "submitted%2Cpartial-filled")
	if err != nil {
		log.Printf("fail sell:%v", err)
	}
	for _, v := range ss.Data {
		sl = append(sl, strconv.Itoa(v.ID))
	}
	log.Printf("fail sell:%v", sl)
	if len(bl) > 50 {
		bl = bl[:50]
	}
	bh, err := h.OrdersCancel(bl)
	if err != nil {
		log.Printf("fail cancel buy:%v", err)
	}
	log.Printf("success  cancel  buy:%v", bh)
	if len(sl) > 50 {
		sl = sl[:50]
	}
	sh, err := h.OrdersCancel(sl)
	if err != nil {
		log.Printf("fail cancel sell:%v", err)
	}
	log.Printf("fail cancel buy:%v", sh)
}

func TestIncrDecimal(t *testing.T) {
	a := 0.00000058
	d := decimal.NewFromFloat(a)
	r, _ := calDecimal(d, 11, 1, "-")
	log.Println(r)
}

func TestGetPrec(t *testing.T) {
	md, err := getHuobiSDK().MarketDepth("pnteth")
	if err != nil {
		log.Fatal("获取井深数据错误", err)
	}
	var s []float64
	for _, v := range md.Tick.Bids {
		s = append(s, v[0])
	}
	res, _ := getPrec(s)
	log.Printf("精度长度是:%v", res)
}

func getHuobiSDK() *robot.HuobiSdk {
	hsdk := &robot.HuobiSdk{
		Host:        "api.huobi.pro",
		AccessKeyID: "2a06b82f-ae8f772e-416b1e23-5aa47",
		SecretKey:   "f4b66904-d6717085-63661938-89805",
	}
	return hsdk
}

func TestContinuousNumPlus(t *testing.T) {
	a := 0.000000999
	d := decimal.NewFromFloat(a)
	r, _ := continuousNums(d, 11, 100, "+")
	for _, v := range r {
		log.Println(v)
	}
}

func TestContinuousNumMinus(t *testing.T) {
	a := 0.0000001
	d := decimal.NewFromFloat(a)
	r, _ := continuousNums(d, 11, 100, "-")
	for _, v := range r {
		log.Println(v)
	}
}

func TestRandNumsPlus(t *testing.T) {
	a := 0.000000999
	d := decimal.NewFromFloat(a)
	r, _ := randNums(d, 11, 100, 2, 5, "+")
	for _, v := range r {
		log.Println(v)
	}
}

func TestRandNumsMinus(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	a := 0.000000560
	d := decimal.NewFromFloat(a)
	r, err := randNums(d, 11, 10, 2, 5, "-")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range r {
		log.Println(v)
	}
}

func TestMergeSlice(t *testing.T) {
	a := []string{"1", "2"}
	b := []string{"3", "4"}
	c := mergeSlice(a, b)
	for _, v := range c {
		log.Println(v)
	}
}

func TestSortStep(t *testing.T) {
	a := StepList{
		robot.Step{
			FlagStart: 1,
			FlagEnd:   5,
			AmountMin: 1,
			AmountMax: 3,
		},
		robot.Step{
			FlagStart: 6,
			FlagEnd:   10,
			AmountMin: 3,
			AmountMax: 5,
		},
		robot.Step{
			FlagStart: 11,
			FlagEnd:   15,
			AmountMin: 5,
			AmountMax: 7,
		},
		robot.Step{
			FlagStart: 16,
			FlagEnd:   20,
			AmountMin: 7,
			AmountMax: 9,
		},
	}
	sort.Sort(a)
	for _, v := range a {
		log.Println(v)
	}

}

func getPlaceConfig() robot.Place {
	pc := robot.Place{
		Symbol: "pnteth",
		Step: []robot.Step{
			robot.Step{
				FlagStart: 1,
				FlagEnd:   5,
				AmountMin: 1,
				AmountMax: 3,
			},
			robot.Step{
				FlagStart: 6,
				FlagEnd:   10,
				AmountMin: 3,
				AmountMax: 5,
			},
			robot.Step{
				FlagStart: 11,
				FlagEnd:   15,
				AmountMin: 5,
				AmountMax: 7,
			},
			robot.Step{
				FlagStart: 16,
				FlagEnd:   20,
				AmountMin: 7,
				AmountMax: 9,
			},
		},
	}
	return pc
}

func TestDecimalEqual(t *testing.T) {
	a := decimal.NewFromFloat(1.01).String()
	b := decimal.NewFromFloat(1.01).String()

	set := make(map[string]struct{})
	set[a] = struct{}{}
	set[b] = struct{}{}
	log.Printf("结果是:%v", set)
}

func TestSortString(t *testing.T) {

	origin := []string{"1", "2", "3"}
	for i := 0; i < len(origin); i++ {
		if origin[i] == "2" {
			origin = append(origin[:i], origin[i+1:]...)
			i--
		}
	}
	log.Println(origin)
}

func TestStrings(t *testing.T) {
	name := [][]int{}
	n1 := []int{1, 2, 3}
	n2 := []int{1, 2, 3}
	name = append(name, n1)
	name = append(name, n2)

	log.Printf("value is :%v", len(name))
}
