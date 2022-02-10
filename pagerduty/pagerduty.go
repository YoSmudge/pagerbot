package pagerduty

import (
	"fmt"
	"gopkg.in/jmcvetta/napping.v3"
	"net/http"
)

type Api struct {
	key      string
	org      string
	http     *napping.Session
	timezone string
}

// Pagerduty API doesn't provide a sane way of checking for auth
// so we just get the schedules at setup time
func New(key string, org string) (*Api, error) {
	a := Api{}
	a.key = key
	a.org = org
	a.timezone = "UTC"

	a.http = &napping.Session{
		Header: &http.Header{
			"Authorization": []string{fmt.Sprintf("Token token=%s", a.key)},
			"User-Agent":    []string{"PagerBot +https://github.com/yosmudge/pagerbot"},
		},
	}

	_, err := a.request("schedules")
	if err != nil {
		return &a, err
	}

	return &a, nil
}
