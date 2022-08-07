package main

import (
	"context"
	"lokli/lcontext"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func GetReverseProxy(ctx context.Context, urlstr string) *httputil.ReverseProxy {
	url.Parse(urlstr)
	return &httputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			outReq.URL.Scheme = "http"
			outReq.URL.Host = urlstr
			outReq.Method = "GET"
			outReq.Header.Set("name", "lokli")
		},

		Transport: &http.Transport{},
	}
}

func Handle1(rsp http.ResponseWriter, req *http.Request) {
	tm := rand.Intn(1000)
	time.Sleep(time.Duration(tm) * time.Millisecond)
	rsp.Write([]byte("handle1"))
}

func Handle2(rsp http.ResponseWriter, req *http.Request) {
	rsp.Write([]byte("handle2"))
}
func main() {
	h := lcontext.NewLokContext("/h1")
	mux := http.NewServeMux()
	mux.HandleFunc("/h1", Handle1)
	mux.HandleFunc("/h2", Handle2)
	mux.Handle("/h", h)

	http.ListenAndServe(":9200", mux)
}
