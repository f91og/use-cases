package main

import (
	"net/http"

	"github.com/golang/glog"
)

type WebhookServer struct {
	server *http.Server
}

func (ws *WebhookServer) validate(w http.ResponseWriter, r *http.Request) {
	// ......
}

func (ws *WebhookServer) mutate(w http.ResponseWriter, r *http.Request) {
	// ......
}

func main() {
	// some logics, maybe parse arguments

	// 启动一个server，处理kube api server过来的validate请求和mutate请求
	whsvr := &WebhookServer{
		server: &http.Server{
			// Addr:      fmt.Sprintf(":%v", parameters.port),
			// TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", whsvr.validate)
	mux.HandleFunc("/mutate", whsvr.mutate)
	mux.HandleFunc("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("running"))
	}))
	whsvr.server.Handler = mux

	// start webhook server in new routine
	if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
		glog.Errorf("Failed to listen and serve webhook server: %v", err)
	}
}
