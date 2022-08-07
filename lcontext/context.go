package lcontext

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lokli/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type LokContext struct {
	ctx           context.Context
	Url           string
	Timout        time.Duration
	req           *http.Request
	overloadValue float64
	limit         rate.Limit
}

func NewLokContext(url string) *LokContext {
	return &LokContext{
		Url:           url,
		overloadValue: 30,
		limit:         rate.NewLimiter(500, 1000),
	}
}

func (lc *LokContext) RateAllow(ctx context.Context) error {
	return lc.limit.Wait(ctx)
}

func (lc *LokContext) Allow() (bool, error) {
	use := utils.GetCpuUsage()
	fmt.Println("cpu usage:", use)
	if use > lc.overloadValue {
		return false, nil
	}
	return true, nil
}

func (lc *LokContext) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	if err := lc.RateAllow(req.Context()); err != nil {
		rsp.Write([]byte(" rate limit exceeded"))
		return
	}
	if ok, err := lc.Allow(); !ok || err != nil {
		rsp.Write([]byte("cpu protect self"))
		return
	}
	lc.req = req
	lc.ctx = req.Context()
	// proxy := httputil.NewSingleHostReverseProxy(u)
	// proxy.Director =
	proxy := &httputil.ReverseProxy{
		Director:      lc.preDirector,
		FlushInterval: 100 * time.Millisecond,
	}
	proxy.ModifyResponse = lc.afterResponse
	proxy.ServeHTTP(rsp, req)
}

func (lc *LokContext) preDirector(out *http.Request) {
	log.Println(lc.req.Host)
	if strings.HasPrefix(lc.Url, "/") {
		lc.Url = "http://" + lc.req.Host + lc.Url
	}
	b, _ := ioutil.ReadAll(lc.req.Body)
	fmt.Println(string(b))
	u, _ := url.Parse(lc.Url)
	out.Header.Set("name", "lokli")
	out.Method = "GET"
	out.Host = lc.Url
	out.URL = u
	out.URL.Scheme = "http"
	out.Body = io.NopCloser(bytes.NewReader(b))
}

func (lc *LokContext) afterResponse(rsp *http.Response) error {
	rsp.Header.Set("rsp", "lok")
	return nil
}
