package server

import (
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"sync"
)

var sessionMap = sync.Map{}

func addToSessionMap(session *MediaSession) {
	sessionMap.Store(session.sessionId, session)
	prom.CreatedSession.Inc()
}

func removeFromSessionMap(session *MediaSession) {
	sessionMap.Delete(session.sessionId)
	prom.CreatedSession.Dec()
}

func (srv *MediaServer) getNextAvailableRtpPort() uint16 {
	return srv.portPool.get()
}

func (srv *MediaServer) reclaimRtpPort(port uint16) {
	if port != 0 {
		srv.portPool.put(port)
	}
}

func (srv *MediaServer) createSession(param *rpc.CreateParam) (session *MediaSession, err error) {
	defer func() {
		if err != nil && session != nil {
			session.finalize()
		}
	}()

	if session, err = newSession(srv, param); err != nil {
		return
	}
	// initialize source/sink list for each session
	// the factory's order is important
	for _, factory := range srv.sourceF {
		src := factory.NewSource(session)
		session.source = append(session.source, src)
	}
	for _, factory := range srv.sinkF {
		sink := factory.NewSink(session)
		session.sink = append(session.sink, sink)
	}

	// connect source/sink into event graph of this session
	// then listen on udp messages
	if err = session.activate(); err != nil {
		return
	}
	prom.AllSession.Inc()
	addToSessionMap(session)
	return session, nil
}
