package collector

import (
	"fmt"

	"dpvs_exporter/lb"

	"github.com/prometheus/client_golang/prometheus"
)

var connInfo map[string]*ConnectionIndicators

// ConnectionIndicators: the indicators of conns
type ConnectionIndicators struct {
	conns    *prometheus.Desc
	inBytes  *prometheus.Desc
	outBytes *prometheus.Desc
	inPkts   *prometheus.Desc
	outPkts  *prometheus.Desc
}

type ConnStatsController struct {
	comm *lb.DpvsAgentComm
}

func NewConnStatsController(agent *lb.DpvsAgentComm) *ConnStatsController {
	return &ConnStatsController{
		comm: agent,
	}
}

// Describe 用于注册 connInfo 中所有的描述符
func (c *ConnStatsController) Describe(ch chan<- *prometheus.Desc) {
	for _, ci := range connInfo {
		// 注册每个指标
		ch <- ci.conns
		ch <- ci.inBytes
		ch <- ci.outBytes
		ch <- ci.inPkts
		ch <- ci.outPkts
	}
}

func (c *ConnStatsController) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[panic] connRateCollector: %v", r)
		}
	}()
	services, err := c.comm.ListVirtualServices()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if services == nil {
		DefaultEmitMissingMetrics(ch, connInfo)
		return
	}
	for _, vss := range services.Items {
		// 使用 GetServerIdentifier 生成唯一的 key
		key := GetServerIdentifier(vss.Addr, vss.Port, vss.Proto)

		// 为每个 service 初始化 ConnectionIndicators
		ci, exists := connInfo[key]
		if exists {
			ch <- prometheus.MustNewConstMetric(ci.conns, prometheus.CounterValue, float64(safeDereferenceInt64(vss.Stats.Conns)), key)
			ch <- prometheus.MustNewConstMetric(ci.inBytes, prometheus.CounterValue, float64(safeDereferenceInt64(vss.Stats.InBytes)), key)
			ch <- prometheus.MustNewConstMetric(ci.outBytes, prometheus.CounterValue, float64(safeDereferenceInt64(vss.Stats.OutBytes)), key)
			ch <- prometheus.MustNewConstMetric(ci.inPkts, prometheus.CounterValue, float64(safeDereferenceInt64(vss.Stats.InPkts)), key)
			ch <- prometheus.MustNewConstMetric(ci.outPkts, prometheus.CounterValue, float64(safeDereferenceInt64(vss.Stats.OutPkts)), key)
		}
		if vss.RSs != nil {
			for _, rs := range vss.RSs.Items {
				rsKey := GetServerIdentifier(rs.Spec.IP, rs.Spec.Port, vss.Proto)
				cii, exists := connInfo[rsKey]
				if exists {
					ch <- prometheus.MustNewConstMetric(cii.conns, prometheus.CounterValue, float64(safeDereferenceInt64(rs.Stats.Conns)), rsKey)
					ch <- prometheus.MustNewConstMetric(cii.inBytes, prometheus.CounterValue, float64(safeDereferenceInt64(rs.Stats.InBytes)), rsKey)
					ch <- prometheus.MustNewConstMetric(cii.outBytes, prometheus.CounterValue, float64(safeDereferenceInt64(rs.Stats.OutBytes)), rsKey)
					ch <- prometheus.MustNewConstMetric(cii.inPkts, prometheus.CounterValue, float64(safeDereferenceInt64(rs.Stats.InPkts)), rsKey)
					ch <- prometheus.MustNewConstMetric(cii.outPkts, prometheus.CounterValue, float64(safeDereferenceInt64(rs.Stats.OutPkts)), rsKey)
				}
			}
		}
	}
}

// InitConnStatsController 初始化连接统计控制器
func InitConnStatsController(services []lb.VirtualServerSpecExpand) {
	// 清空 connInfo
	connInfo = make(map[string]*ConnectionIndicators)

	// 遍历服务列表，初始化每个服务的连接指标
	for _, vss := range services {
		// 生成唯一的服务器标识符
		key := GetServerIdentifier(vss.Addr, vss.Port, vss.Proto)
		// 创建 ConnectionIndicators
		connInfo[key] = &ConnectionIndicators{
			conns: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "conn", key+"_conns"),
				"VIP connections",
				[]string{"conns"},
				nil,
			),
			inBytes: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "conn", key+"_in_bytes"),
				"Incoming bytes for VIP",
				[]string{"inBytes"},
				nil,
			),
			outBytes: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "conn", key+"_out_bytes"),
				"Outgoing bytes for VIP",
				[]string{"outBytes"},
				nil,
			),
			inPkts: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "conn", key+"_in_pkts"),
				"Incoming packets for VIP",
				[]string{"inPkts"},
				nil,
			),
			outPkts: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "conn", key+"_out_pkts"),
				"Outgoing packets for VIP",
				[]string{"outPkts"},
				nil,
			),
		}
		if vss.RSs != nil {
			for _, rs := range vss.RSs.Items {
				rsKey := GetServerIdentifier(rs.Spec.IP, rs.Spec.Port, vss.Proto)
				connInfo[rsKey] = &ConnectionIndicators{
					conns: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "conn", rsKey+"_conns"),
						"RS connections",
						[]string{"conns"},
						nil,
					),
					inBytes: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "conn", rsKey+"_in_bytes"),
						"Incoming bytes for RS",
						[]string{"inBytes"},
						nil,
					),
					outBytes: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "conn", rsKey+"_out_bytes"),
						"Outgoing bytes for RS",
						[]string{"outBytes"},
						nil,
					),
					inPkts: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "conn", rsKey+"_in_pkts"),
						"Incoming packets for RS",
						[]string{"inPkts"},
						nil,
					),
					outPkts: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "conn", rsKey+"_out_pkts"),
						"Outgoing packets for RS",
						[]string{"outPkts"},
						nil,
					),
				}
			}
		}
	}
}

// 拼接 addr, port 和 proto 成 ip:port:proto 格式的函数
func GetServerIdentifier(addr *string, port *int64, proto *int64) string {
	// 如果 addr 是 nil，使用默认的 IP 地址
	if addr == nil {
		addr = new(string)
		*addr = "0.0.0.0" // 默认值
	}

	// 如果 port 是 nil，使用默认的端口
	if port == nil {
		port = new(int64)
		*port = 80 // 默认值
	}

	// 如果 proto 是 nil，使用默认的协议（假设为 TCP，即 6）
	if proto == nil {
		proto = new(int64)
		*proto = 6 // 默认值
	}

	// 根据 proto 值判断是 TCP 还是 UDP
	protoStr := "TCP" // 默认是 TCP
	if *proto == 17 {
		protoStr = "UDP"
	}

	// 拼接结果并返回
	return fmt.Sprintf("%s:%d:%s", *addr, *port, protoStr)
}
func DefaultEmitMissingMetrics(ch chan<- prometheus.Metric, connInfo map[string]*ConnectionIndicators) {
	for key, ci := range connInfo {
		ch <- prometheus.MustNewConstMetric(ci.conns, prometheus.GaugeValue, 0, key)
		ch <- prometheus.MustNewConstMetric(ci.inBytes, prometheus.GaugeValue, 0, key)
		ch <- prometheus.MustNewConstMetric(ci.outBytes, prometheus.GaugeValue, 0, key)
		ch <- prometheus.MustNewConstMetric(ci.inPkts, prometheus.GaugeValue, 0, key)
		ch <- prometheus.MustNewConstMetric(ci.outPkts, prometheus.GaugeValue, 0, key)
	}
}
