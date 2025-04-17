package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	NodeSystemEventException = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "node_sys_event_exception",
		Help: "exceptions when handling system event",
	})
	NodeUserEventException = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "node_user_event_exception",
		Help: "Exceptions when handling user event",
	})
	NodeGraphNodes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_graph_nodes",
		Help: "Node number in all graph",
	})
	NodeGraphLinks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_graph_links",
		Help: "Link number in all graph",
	})
	RtpCreatedSession = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rtp_created_session",
		Help: "Created session",
	})
	RtpStartedSession = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rtp_started_session",
		Help: "Created as well as started session",
	})
	RtpAllSession = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rtp_all_session",
		Help: "Total created sessions since start-up",
	})
	RtpAbnormalSession = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rtp_abnormal_session",
		Help: "Abnormal exited rtp session(rtp/rtcp loops not fully exited)",
	})
	RtpSessionGoroutine = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rtp_session_goroutine",
		Help: "goroutine for session(send,recv,recv_ctrl)",
	}, []string{"type"})
	RtpUsedPortPair = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rtp_used_port_pair",
		Help: "Port pairs(rtp/rtcp) allocated",
	})
	GrpcSessionAction = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_session_action",
		Help: "Executed action on session",
	}, []string{"cmd", "type"})
)

func InitCollector() {
	var cs = []prometheus.Collector{
		NodeSystemEventException,
		NodeUserEventException,
		NodeGraphNodes,
		NodeGraphLinks,

		RtpCreatedSession,
		RtpStartedSession,
		RtpAllSession,
		RtpAbnormalSession,
		GrpcSessionAction,
		RtpSessionGoroutine,
		RtpUsedPortPair,
	}
	for _, c := range cs {
		prometheus.MustRegister(c)
	}
}
