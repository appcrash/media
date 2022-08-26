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
	CreatedSession = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "created_session",
		Help: "Created session",
	})
	StartedSession = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "started_session",
		Help: "Created as well as started session",
	})
	AllSession = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "all_session",
		Help: "Total created sessions since start-up",
	})

	SessionAction = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "session_action",
		Help: "Executed action on session",
	}, []string{"cmd", "type"})
	SessionGoroutine = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "session_goroutine",
		Help: "goroutine for session(send,recv,recv_ctrl)",
	}, []string{"type"})
	UsedPortPair = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "used_port_pair",
		Help: "Port pairs allocated",
	})
)

func InitCollector() {
	var cs = []prometheus.Collector{
		NodeSystemEventException,
		NodeUserEventException,
		NodeGraphNodes,
		NodeGraphLinks,

		CreatedSession,
		StartedSession,
		AllSession,
		SessionAction,
		SessionGoroutine,
		UsedPortPair,
	}
	for _, c := range cs {
		prometheus.MustRegister(c)
	}
}
