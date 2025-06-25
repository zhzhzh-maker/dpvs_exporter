package lb

import (
	"net"

	"dpvs_exporter/utils"
)

type Checker uint16

const (
	CheckerNone Checker = iota
	CheckerTCP
	CheckerUDP
	CheckerPING
	CheckerUDPPing
)

type RealServer struct {
	IP        net.IP
	Port      uint16
	Weight    uint16
	Inhibited bool
}

type VirtualService struct {
	Id       string
	Checker  Checker
	Protocol utils.IPProto
	Port     uint16
	IP       net.IP
	RSs      []RealServer
}

type Comm interface {
	ListVirtualServices() ([]VirtualService, error)
	UpdateByChecker(targets []VirtualService) error
}

func (checker Checker) String() string {
	switch checker {
	case CheckerNone:
		return "checker_none"
	case CheckerTCP:
		return "checker_tcp"
	case CheckerUDP:
		return "checker_udp"
	case CheckerPING:
		return "checker_ping"
	}
	return "checker_unknown"
}

// NicDeviceSpecList
type NICStatsResponse struct {
	Items []NICDeviceSpec `json:"Items,omitempty"`
}

// NicDeviceSpec
type NICDeviceSpec struct {
	Detail *NICDeviceDetail `json:"detail,omitempty"`
	Stats  *NICDeviceStats  `json:"stats,omitempty"`
}

// NicDeviceDetail
type NICDeviceDetail struct {
	Addr     *string  `json:"addr,omitempty"`
	Autoneg  *Autoneg `json:"autoneg,omitempty"`
	Duplex   *Duplex  `json:"duplex,omitempty"`
	Flags    *int64   `json:"Flags,omitempty"`
	ID       *int64   `json:"ID,omitempty"`
	MTU      *int64   `json:"MTU,omitempty"`
	Name     *string  `json:"name,omitempty"`
	NRxQ     *int64   `json:"nRxQ,omitempty"`
	NTxQ     *int64   `json:"nTxQ,omitempty"`
	SocketID *int64   `json:"socketID,omitempty"`
	Speed    *int64   `json:"speed,omitempty"`
	Status   *Status  `json:"status,omitempty"`
}

// NicDeviceStats
type NICDeviceStats struct {
	BufAvail    *int64  `json:"bufAvail,omitempty"`
	BufInuse    *int64  `json:"bufInuse,omitempty"`
	ErrorBytesQ []int64 `json:"errorBytesQ,omitempty"`
	ID          *int64  `json:"id,omitempty"`
	InBytes     *int64  `json:"inBytes,omitempty"`
	InBytesQ    []int64 `json:"inBytesQ,omitempty"`
	InErrors    *int64  `json:"inErrors,omitempty"`
	InMissed    *int64  `json:"inMissed,omitempty"`
	InPkts      *int64  `json:"inPkts,omitempty"`
	InPktsQ     []int64 `json:"inPktsQ,omitempty"`
	OutBytes    *int64  `json:"outBytes,omitempty"`
	OutBytesQ   []int64 `json:"outBytesQ,omitempty"`
	OutErrors   *int64  `json:"outErrors,omitempty"`
	OutPkts     *int64  `json:"outPkts,omitempty"`
	OutPktsQ    []int64 `json:"outPktsQ,omitempty"`
	RxNoMbuf    *int64  `json:"rxNoMbuf,omitempty"`
	Name        *string `json:"-"` // 忽略该字段
}

type Autoneg string

const (
	AutoNego  Autoneg = "auto-nego"
	FixedNego Autoneg = "fixed-nego"
)

type Duplex string

const (
	FullDuplex Duplex = "full-duplex"
	HalfDuplex Duplex = "half-duplex"
)

type Status string

const (
	Down Status = "DOWN"
	Up   Status = "UP"
)

// Request
//
// VirtualServerList
type VsResponse struct {
	Items []VirtualServerSpecExpand `json:"Items,omitempty"`
}

// VirtualServerSpecExpand
type VirtualServerSpecExpand struct {
	Addr            *string          `json:"Addr,omitempty"`
	AF              *int64           `json:"Af,omitempty"`
	Bps             *int64           `json:"Bps,omitempty"`
	ConnTimeout     *int64           `json:"ConnTimeout,omitempty"`
	DestCheck       []DestCheckSpec  `json:"DestCheck,omitempty"`
	ExpireQuiescent *ExpireQuiescent `json:"ExpireQuiescent,omitempty"`
	Flags           *string          `json:"Flags,omitempty"`
	Fwmark          *int64           `json:"Fwmark,omitempty"`
	LimitProportion *int64           `json:"LimitProportion,omitempty"`
	Match           *MatchSpec       `json:"Match,omitempty"`
	Netmask         *int64           `json:"Netmask,omitempty"`
	Port            *int64           `json:"Port,omitempty"`
	Proto           *int64           `json:"Proto,omitempty"`
	// 0  (0x00): disable
	// 1  (0x01): v1
	// 2  (0x02): v2
	// 17 (0x11): v1-insecure
	// 18 (0x12): v2-insecure
	ProxyProto *int64                `json:"ProxyProto,omitempty"`
	RSs        *RealServerExpandList `json:"RSs,omitempty"`
	SchedName  *SchedName            `json:"SchedName,omitempty"`
	Stats      *ServerStats          `json:"Stats,omitempty"`
	SynProxy   *ExpireQuiescent      `json:"SynProxy,omitempty"`
	Timeout    *int64                `json:"Timeout,omitempty"`
}

// MatchSpec
type MatchSpec struct {
	Dest      *AddrRange `json:"Dest,omitempty"`
	InIfName  *string    `json:"InIfName,omitempty"`
	OutIfName *string    `json:"OutIfName,omitempty"`
	Src       *AddrRange `json:"Src,omitempty"`
}

// AddrRange
type AddrRange struct {
	End   *string `json:"End,omitempty"`
	Start *string `json:"Start,omitempty"`
}

// RealServerExpandList
type RealServerExpandList struct {
	Items []RealServerSpecExpand `json:"Items,omitempty"`
}

// RealServerSpecExpand
type RealServerSpecExpand struct {
	Spec  *RealServerSpecTiny `json:"Spec,omitempty"`
	Stats *ServerStats        `json:"Stats,omitempty"`
}

// RealServerSpecTiny
type RealServerSpecTiny struct {
	Inhibited  *bool   `json:"inhibited,omitempty"`
	IP         *string `json:"ip,omitempty"`
	Mode       *Mode   `json:"mode,omitempty"`
	Overloaded *bool   `json:"overloaded,omitempty"`
	Port       *int64  `json:"port,omitempty"`
	Weight     *int64  `json:"weight,omitempty"`
}

// ServerStats
type ServerStats struct {
	Conns    *int64 `json:"Conns,omitempty"`
	CPS      *int64 `json:"Cps,omitempty"`
	InBps    *int64 `json:"InBps,omitempty"`
	InBytes  *int64 `json:"InBytes,omitempty"`
	InPkts   *int64 `json:"InPkts,omitempty"`
	InPps    *int64 `json:"InPps,omitempty"`
	OutBps   *int64 `json:"OutBps,omitempty"`
	OutBytes *int64 `json:"OutBytes,omitempty"`
	OutPkts  *int64 `json:"OutPkts,omitempty"`
	OutPps   *int64 `json:"OutPps,omitempty"`
}

// DestCheckSpec
type DestCheckSpec string

const (
	Passive DestCheckSpec = "passive"
	Ping    DestCheckSpec = "ping"
	TCP     DestCheckSpec = "tcp"
	UDP     DestCheckSpec = "udp"
)

type ExpireQuiescent string

const (
	False ExpireQuiescent = "false"
	True  ExpireQuiescent = "true"
)

type Mode string

const (
	DR     Mode = "DR"
	Fnat   Mode = "FNAT"
	Nat    Mode = "NAT"
	Snat   Mode = "SNAT"
	Tunnel Mode = "TUNNEL"
)

type SchedName string

const (
	Conhash SchedName = "conhash"
	Rr      SchedName = "rr"
	Wlc     SchedName = "wlc"
	Wrr     SchedName = "wrr"
)
