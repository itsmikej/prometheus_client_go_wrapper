package main

import (
	"math/rand"
	"time"

	prometheusWrapper "github.com/itsmikej/prometheus_client_go_wrapper"
	"github.com/prometheus/client_golang/prometheus"
)

var ins *prometheusWrapper.PrometheusWrapper

func main() {
	var err error
	ins, err = prometheusWrapper.NewPrometheusWrapper(&prometheusWrapper.Config{
		App:       "test",
		Idc:       "beijing",
		LogMethod: []string{"GET", "POST"},
		LogApi:    []string{"/foo", "/bar"},
		Buckets:   prometheus.LinearBuckets(10, 10, 20),
		Service:   struct{ ListenPort int }{ListenPort: 9000},
	})
	if err != nil {
		panic(err)
	}

	go autoLog()

	t := time.NewTicker(time.Second * 5)
	for {
		<-t.C
		// 业务中自定义的log
		ins.RequestLog("backend", "/baz", "GET", "200")
		ins.RequestLog("backend", "/baz", "GET", "500")
		ins.RcvdBytesLog("backend", "/baz", "GET", "200", 100)
		ins.SendBytesLog("backend", "/baz", "GET", "200", 3000)
		ins.LatencyLog("backend", "/baz", "GET", float64(rand.Intn(200)))
		ins.StateLog("reading", 500)
		ins.ExceptionLog("mysql", "timeout")
		ins.ExceptionLog("mysql", "panic")
	}
}

func autoLog() {
	// 模拟请求结束后记录日志，这块逻辑应该写到请求流程最后或者中间件中
	t := time.NewTicker(time.Second * 1)
	for {
		<-t.C
		ins.Log("/foo", "GET", "200", 5000, 200, float64(rand.Intn(200)))
		ins.Log("/baz", "GET", "200", 5000, 200, float64(rand.Intn(200))) // 将会被过滤
	}
}