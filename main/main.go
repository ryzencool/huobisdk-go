package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/elgs/cron"
	"github.com/ryzencool/huobisdk-go"
	"github.com/ryzencool/huobisdk-go/service"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	hbsk := getHuobiSdk()
	mc := robot.GetMatchConfig()
	c := cron.New()
	cjob1 := 111111
	cjob2 := 111111
	cjob3 := 111111
	cjob4 := 111111
	for _, v := range mc.Match {
		c.AddFunc("0 */5 * * * ?", func() {

			if cjob1 != 111111 {
				c.RemoveFunc(cjob1)
			}
			if cjob2 != 111111 {
				c.RemoveFunc(cjob2)
			}
			if cjob3 != 111111 {
				c.RemoveFunc(cjob3)
			}
			if cjob4 != 111111 {
				c.RemoveFunc(cjob4)
			}
			nh := time.Now().Hour()

			if nh >= 1 && nh < 7 {
				n1 := rand.Intn(len(v.Strategy.Night.Normal))
				n2 := rand.Intn(len(v.Strategy.Night.Special))
				nn := v.Strategy.Night.Normal[n1]
				ns := v.Strategy.Night.Special[n2]
				cjob1, _ = c.AddFunc("*/"+strconv.Itoa(nn.Interval)+" * * * * ?", func() {

					service.AutoMatch(hbsk, v.Symbol, int64(v.Limit), nn.AmountMin, nn.AmountMax)
				})
				cjob2, _ = c.AddFunc("*/"+strconv.Itoa(ns.Interval)+" * * * * ?", func() {

					time.Sleep(time.Duration(rand.Intn(ns.Interval-10)) * time.Second)
					service.AutoMatch(hbsk, v.Symbol, int64(v.Limit), ns.AmountMin, ns.AmountMax)
				})
			} else {
				// 白天
				n3 := rand.Intn(len(v.Strategy.Day.Normal))
				n4 := rand.Intn(len(v.Strategy.Day.Special))
				dn := v.Strategy.Day.Normal[n3]
				ds := v.Strategy.Day.Special[n4]
				cjob3, _ = c.AddFunc("*/"+strconv.Itoa(dn.Interval)+" * * * * ?", func() {

					service.AutoMatch(hbsk, v.Symbol, int64(v.Limit), dn.AmountMin, dn.AmountMax)
				})
				cjob4, _ = c.AddFunc("*/"+strconv.Itoa(ds.Interval)+" * * * * ?", func() {

					time.Sleep(time.Duration(rand.Intn((ds.Interval - 10))) * time.Second)
					service.AutoMatch(hbsk, v.Symbol, int64(v.Limit), ds.AmountMin, ds.AmountMax)
				})
			}
		})
	}

	c.Start()
	for _, v := range mc.Place {
		go func() {
			for {
				service.AutoPlace(hbsk, v)
			}
		}()
		for i := 0; i < 3; i++ {
			go func() {
				for {
					service.CancelOver30(hbsk, v.Symbol)
				}
			}()
		}
	}

	select {}
}

func getHuobiSdk() *robot.HuobiSdk {
	var config = robot.GetMatchConfig()
	h := &robot.HuobiSdk{
		Host:        config.Host,
		AccessKeyID: config.AccessID,
		SecretKey:   config.SecretKey,
	}
	return h
}
