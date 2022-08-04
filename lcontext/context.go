package lcontext

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LokContext struct {
	ctx    context.Context
	Url    string
	Timout time.Duration
	req    *http.Request
}

func NewLokContext(url string) *LokContext {
	return &LokContext{
		Url: url,
	}
}

func (lc *LokContext) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
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
