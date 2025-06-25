// package main

// import (
// 	"dpvs_exporter/collector"
// 	"log"
// 	"log/slog"
// 	"net/http"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// )

// func main() {
// 	handler := slog.NewTextHandler(log.Writer(), nil)
// 	logger := slog.New(handler)
// 	// dpvsCollector := collector.NewDpvsCollector(logger)
// 	netCollector := collector.NewNetDevCollector(logger)
// 	// 将 DpvsCollector 注册到 Prometheus
// 	prometheus.MustRegister(netCollector)

// 	// 设置 HTTP handler
// 	http.Handle("/metrics", promhttp.Handler())

//		// 启动 HTTP 服务
//		log.Println("Starting server on :9101")
//		if err := http.ListenAndServe(":9101", nil); err != nil {
//			log.Fatalf("Error starting server: %v", err)
//		}
//	}
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"dpvs_exporter/collector"
	"dpvs_exporter/lb"
)

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9101", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()
	agent := lb.NewDpvsAgentComm("")
	nicName, err := agent.ListNicName()
	if err != nil || nicName == nil {
		fmt.Println(err.Error())
		return
	}
	serverInfo, err := agent.ListVirtualServices()
	if err != nil || serverInfo == nil {
		return
	}
	collector.InitConnStatsController(serverInfo.Items)
	collector.InitNicCollector(nicName)
	dpvs := collector.NewDpvs(*agent)
	prometheus.MustRegister(dpvs)

	http.Handle(*metricsPath, promhttp.Handler())
	log.Printf("Starting dpvs_exporter on %s%s\n", *listenAddress, *metricsPath)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
