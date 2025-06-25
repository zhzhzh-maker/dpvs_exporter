package collector

import (
	"dpvs_exporter/lb"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "dpvs"
	subsystem = "nic"
	clientID  = "dpvs_exporter"
)

var labels = prometheus.Labels{"source": "dpvs-agent"}

type Snap struct {
	buffAvail *prometheus.Desc
	buffInUse *prometheus.Desc
	inBytes   *prometheus.Desc
	outBytes  *prometheus.Desc
	inPkts    *prometheus.Desc
	outPkts   *prometheus.Desc
	inErrors  *prometheus.Desc
}

var nics map[string]*Snap

type NicRateCollector struct {
	comm *lb.DpvsAgentComm
}

func NewNicRateCollector(comm *lb.DpvsAgentComm) *NicRateCollector {
	return &NicRateCollector{
		comm: comm,
	}
}

func (c *NicRateCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, nic := range nics {
		describeSnap(ch, nic)
	}
}

func describeSnap(ch chan<- *prometheus.Desc, snap *Snap) {
	if snap == nil {
		return
	}
	if snap.buffAvail != nil {
		ch <- snap.buffAvail
	}
	if snap.buffInUse != nil {
		ch <- snap.buffInUse
	}
	if snap.inBytes != nil {
		ch <- snap.inBytes
	}
	if snap.outBytes != nil {
		ch <- snap.outBytes
	}
	if snap.inPkts != nil {
		ch <- snap.inPkts
	}
	if snap.outPkts != nil {
		ch <- snap.outPkts
	}
	if snap.inErrors != nil {
		ch <- snap.inErrors
	}
}

// 模拟生成网卡数据的结构体
type NicStats struct {
	Name      string
	BuffAvail int64
	BuffInUse int64
	InBytes   int64
	OutBytes  int64
	InPkts    int64
	OutPkts   int64
	InErrors  int64
}

// 模拟生成一些随机网卡数据
func (c *NicRateCollector) generateRandomNicStats() []NicStats {
	nicStats, err := c.comm.ListNicStats()
	stats := make([]NicStats, 0, len(nicStats))
	if err != nil {
		return stats
	}
	for _, nic := range nicStats {
		// 确保字段是非 nil 才进行解引用
		stats = append(stats, NicStats{
			Name:      safeDereference(nic.Name),
			BuffAvail: safeDereferenceInt64(nic.BufAvail),
			BuffInUse: safeDereferenceInt64(nic.BufInuse),
			InBytes:   safeDereferenceInt64(nic.InBytes),
			OutBytes:  safeDereferenceInt64(nic.OutBytes),
			InPkts:    safeDereferenceInt64(nic.InPkts),
			OutPkts:   safeDereferenceInt64(nic.OutPkts),
			InErrors:  safeDereferenceInt64(nic.InErrors),
		})
	}

	return stats
}

func (c *NicRateCollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[panic] NicRateCollector: %v", r)
		}
	}()
	// 模拟调用comm.ListNicStats()改成用模拟数据
	nicStats := c.generateRandomNicStats()

	for _, stat := range nicStats {
		nic, exists := nics[stat.Name]
		if exists {
			// ch <- prometheus.MustNewConstMetric(nic.buffAvail, prometheus.GaugeValue, float64(stat.BuffAvail), stat.Name)
			// ch <- prometheus.MustNewConstMetric(nic.buffInUse, prometheus.GaugeValue, float64(stat.BuffInUse), stat.Name)
			ch <- prometheus.MustNewConstMetric(nic.inBytes, prometheus.CounterValue, float64(stat.InBytes), stat.Name)
			ch <- prometheus.MustNewConstMetric(nic.inPkts, prometheus.CounterValue, float64(stat.InPkts), stat.Name)
			ch <- prometheus.MustNewConstMetric(nic.outBytes, prometheus.CounterValue, float64(stat.OutBytes), stat.Name)
			ch <- prometheus.MustNewConstMetric(nic.outPkts, prometheus.CounterValue, float64(stat.OutPkts), stat.Name)
			ch <- prometheus.MustNewConstMetric(nic.inErrors, prometheus.CounterValue, float64(stat.InErrors), stat.Name)
		}

	}
}

func InitNicCollector(nicName []string) {
	nics = make(map[string]*Snap, 0)
	for _, name := range nicName {
		value := &Snap{
			buffAvail: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_buff_available"),
				"Available buffer count for incoming packets.",
				[]string{"nic"}, labels,
			),
			buffInUse: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_buff_inuse"),
				"In-use buffer count for incoming packets.",
				[]string{"nic"}, labels,
			),
			inBytes: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_in_bytes"),
				"Bytes received.",
				[]string{"nic"}, labels,
			),
			inPkts: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_in_packets"),
				"Packets received.",
				[]string{"nic"}, labels,
			),
			outBytes: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_out_bytes"),
				"Bytes received.",
				[]string{"nic"}, labels,
			),
			outPkts: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_out_packets"),
				"Packets received.",
				[]string{"nic"}, labels,
			),
			inErrors: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name+"_in_errors"),
				"Receive errors.",
				[]string{"nic"}, labels,
			),
		}
		nics[name] = value
	}
}

// 安全的解引用字符串（如果是 nil 返回空字符串）
func safeDereference(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// 安全的解引用 int64（如果是 nil 返回 0）
func safeDereferenceInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
