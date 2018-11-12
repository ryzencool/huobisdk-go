package service

import (
	"log"
	"testing"

	robot "github.com/ryzencool/huobisdk-go"
)

func TestCheckPrice(t *testing.T) {
	log.Println(checkNumNear("0.0000019", "0.000002"))
}

func TestFindRandPrice(t *testing.T) {
	for i := 0; i < 100; i++ {
		log.Println(findRand(100.40, 100.44))
	}
}

func getHuobiSdk() *robot.HuobiSdk {
	h := &robot.HuobiSdk{
		Host:        "api.huobi.pro",
		AccessKeyID: "2a06b82f-ae8f772e-416b1e23-5aa47",
		SecretKey:   "f4b66904-d6717085-63661938-89805",
	}
	return h
}

// func getRobotMatch() robot.Match {
// 	return robot.Match{
// 		Symbol:    "pnteth",
// 		Limit:     10000,
// 		AmountMax: 100,
// 		AmountMin: 1,
// 	}
// }
