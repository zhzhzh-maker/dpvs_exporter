package collector

import (
	"dpvs_exporter/lb"

	"github.com/prometheus/client_golang/prometheus"
)

type Dpvs struct {
	conn *ConnStatsController
	nic  *NicRateCollector
}

func NewDpvs(agent lb.DpvsAgentComm) *Dpvs {
	return &Dpvs{
		conn: NewConnStatsController(&agent),
		nic:  NewNicRateCollector(&agent),
	}
}

func (c *Dpvs) Collect(ch chan<- prometheus.Metric) {
	c.conn.Collect(ch)
	c.nic.Collect(ch)
}

func (c *Dpvs) Describe(ch chan<- *prometheus.Desc) {
	c.conn.Describe(ch)
	c.nic.Describe(ch)
}
