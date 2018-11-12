package robot

import (
	"log"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	log.Println(time.Now().UTC().String())
	var params = make(map[string]string)
	Signature("GET", "api.hadax.com", "/v1/account/accounts", params, "f4b66904-d6717085-63661938-89805", "2a06b82f-ae8f772e-416b1e23-5aa47")
	res, err := Get("https://api.hadax.com/v1/account/accounts", params)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(res)
	}
}
