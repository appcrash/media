package server

import (
	"context"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/utils"
	"github.com/prometheus/client_golang/prometheus"
	"runtime/debug"
)

// receive rtcp packet
func (s *MediaSession) receiveRtcpLoop(ctx context.Context) {
	rtcpReceiver := s.rtpSession.CreateCtrlEventChan()
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "recv_ctrl"})
	gauge.Inc()

	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop ctrl recv", s.GetSessionId())
		s.doneC <- "done"
	}()
	cancelC := ctx.Done()
	for {
		select {
		case eventArray, more := <-rtcpReceiver:
			if !more {
				// RTP stack closed rtcp channel, just return
				return
			}
			for _, evt := range eventArray {
				if evt.EventType == rtp.RtcpBye {
					// peer send bye, notify data send/receive loop to stop
					logger.Debugf("session: %v rtp peer says bye", s.sessionId)
					go s.Stop() // CAVEAT: don't call Stop() in this goroutine directly
					return
				}
			}
		case <-cancelC:
			return
		}
	}
}

func (s *MediaSession) receiveRtpLoop(ctx context.Context) {
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "recv"})
	gauge.Inc()
	// Create and store the data receive channel.
	defer func() {
		if r := recover(); r != nil {
			logger.Errorln("receiveRtpLoop panic(recovered)")
			debug.PrintStack()
		}
	}()

	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop local receive", s.GetSessionId())
		if s.handleC != nil {
			// notify packet handler
			close(s.handleC)
			s.handleC = nil
		}
		s.doneC <- "done"
	}()

	if s.handleC == nil {
		logger.Infof("session:%v has no rtp handling channel, stop local receive early", s.sessionId)
		return
	}

	s.watchdog.reportLoopInfo(receiveLoop)
	rtpSession := s.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()
	cancelC := ctx.Done()
	var nbPacket int
	for {
		select {
		case rp, more := <-dataReceiver:
			if !more {
				// RTP stack closed this channel, so stop receiving anymore
				return
			}

			// nonblock push received data to handler
			pl := utils.NewPacketListFromRtpPacket(rp)
			select {
			case s.handleC <- pl:
			default:
			}
			nbPacket++
			if nbPacket > ReportInfoPacketInterval {
				nbPacket = 0
				s.watchdog.reportLoopInfo(receiveLoop)
			}

			// don't free packet, let it be GCed, as GoRTP will reuse this packet along with its buffer
			// which may be hold by other packet-list objects
			// rp.FreePacket()
		case <-cancelC:
			return
		}
	}
}

func (s *MediaSession) sendRtpLoop(ctx context.Context) {
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "send"})
	gauge.Inc()

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendRtpLoop panic %v", r)
			debug.PrintStack()
		}
	}()
	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop local send", s.GetSessionId())
		s.doneC <- "done"
	}()

	if s.pullC == nil {
		logger.Infof("session:%v has no rtp pulling channel, stop local send early", s.sessionId)
	}

	var nbPacket int
	cancelC := ctx.Done()
	for {
		select {
		// pump data out from graph
		case packetList, more := <-s.pullC:
			if !more {
				return
			}

			if s.rtpSession == nil {
				return
			}
			if packetList == nil {
				continue
			}

			// send all packets based on RtpPacketList
			// for video, a frame can have more than one packet with same timestamp
			packetList.Iterate(func(p *utils.RtpPacketList) {
				payload, _, pts, mark := p.Payload, p.PayloadType, p.Pts, p.Marker
				if payload != nil {
					packet := s.rtpSession.NewDataPacket(pts)
					packet.SetMarker(mark)
					packet.SetPayload(payload)
					//maybe update pt by sip/sdp after create graph
					packet.SetPayloadType(s.avPayloadNumber)
					if _, err := s.rtpSession.WriteData(packet); err != nil {
						s.watchdog.reportLoopError(sendLoop, err)
					}
				}
				nbPacket++
			})
			if nbPacket > ReportInfoPacketInterval {
				nbPacket = 0
				s.watchdog.reportLoopInfo(sendLoop)
			}
		case <-cancelC:
			return
		}
	}

}
