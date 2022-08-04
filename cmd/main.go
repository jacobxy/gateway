package main

import (
	"context"
	"lokli/lcontext"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	rsp.Write([]byte("handle1"))
}

func Handle2(rsp http.ResponseWriter, req *http.Request) {
	rsp.Write([]byte("handle2"))
}
func main() {
	h := lcontext.NewLokContext("http://127.0.0.1:8080/h1")
	mux := http.NewServeMux()
	mux.HandleFunc("/h1", Handle1)
	mux.HandleFunc("/h2", Handle2)
	mux.Handle("/h", h)

	http.ListenAndServe(":8080", mux)
}
