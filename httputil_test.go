package robot

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	now := time.Now()
	res, err := Get("https://api.huobi.pro/market/depth", map[string]string{"symbol": "your symbol", "type": "step0"})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	fmt.Println(time.Since(now))

}

func TestPost(t *testing.T) {
	res, err := Post("http://localhost:8081/test/post", map[string]interface{}{"name": "zmy", "age": "11"})
	if err != nil {
		log.Fatal("fail post something")
	} else {
		log.Panicln("post success:", res)
	}
}
