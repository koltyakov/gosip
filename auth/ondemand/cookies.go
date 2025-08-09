package ondemand

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

var (
	cookiesCheck = []string{"FedAuth", "EdgeAccessCookie"}
)

// Cookies type
type Cookies []Cookie

// Cookie struct
type Cookie struct {
	Domain  string  `json:"domain"`
	Name    string  `json:"name"`
	Value   string  `json:"value"`
	Expires float64 `json:"expires"`
}

func (cookies *Cookies) toString() string {
	res := ""
	for _, cookie := range *cookies {
		if len(res) > 0 {
			res += "; "
		}
		res += fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}
	return res
}

func (cookies *Cookies) isExpired() bool {
	for _, check := range cookiesCheck {
		for _, cookie := range *cookies {
			if cookie.Name == check {
				if cookie.isExpired() {
					return true
				}
			}
		}
	}
	if len(*cookies) == 0 {
		return true
	}
	return false
}

func (cookies *Cookies) getExpire() int64 {
	var exp int64 = -1
	for _, check := range cookiesCheck {
		for _, cookie := range *cookies {
			if cookie.Name == check {
				e := cookie.getExpire()
				if exp == -1 || e < exp {
					exp = e
				}
			}
		}
	}
	return exp
}

func (cookie *Cookie) isExpired() bool {
	if cookie.Expires == -1 {
		return false
	}
	sec, dec := math.Modf(cookie.Expires)
	expireTime := time.Unix(int64(sec), int64(dec*(1e9)))
	if time.Now().Add(time.Minute).Before(expireTime) {
		return false
	}
	return true
}

func (cookie *Cookie) getExpire() int64 {
	sec, dec := math.Modf(cookie.Expires)
	expireTime := time.Unix(int64(sec), int64(dec*(1e9)))
	return expireTime.Unix()
}

func (cookie *Cookie) toMap() map[string]interface{} {
	res := map[string]interface{}{}
	raw, _ := json.Marshal(cookie)
	_ = json.Unmarshal(raw, &res)
	return res
}
