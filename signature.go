package robot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
	"time"
)

const signatureVersion = "2"

const signatureMethod = "HmacSHA256"

// Signature is sign with method, and so on
func Signature(method string, host string, webURL string, params map[string]string, secret string, accessKeyID string) {
	params["AccessKeyId"] = accessKeyID
	params["SignatureVersion"] = signatureVersion
	params["SignatureMethod"] = signatureMethod
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
	origin := method + "\n" + host + "\n" + webURL + "\n" + p
	params["Signature"] = url.QueryEscape(computeHmac256(origin, secret))
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
