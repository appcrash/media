package server

import (
	"context"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/codec"
	"github.com/appcrash/media/server/prom"
	"github.com/prometheus/client_golang/prometheus"
	"runtime/debug"
	"time"
)

// receive rtcp packet
func (s *MediaSession) receiveCtrlLoop(ctx context.Context) {
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
					logger.Debugln("rtp peer says bye")
					go s.Stop() // CAVEAT: don't call Stop() in this goroutine
					return
				}
			}
		case <-cancelC:
			return
		}
	}
}

func (s *MediaSession) receivePacketLoop(ctx context.Context) {
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "recv"})
	gauge.Inc()
	// Create and store the data receive channel.
	defer func() {
		if r := recover(); r != nil {
			logger.Fatalln("receivePacketLoop panic(recovered)")
			debug.PrintStack()
		}
	}()

	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop local receive", s.GetSessionId())
		s.doneC <- "done"
	}()

	rtpSession := s.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()
	cancelC := ctx.Done()
	for {
		select {
		case rp, more := <-dataReceiver:
			if !more {
				// RTP stack closed this channel, so stop receiving anymore
				return
			}

			// push received data to all sinks, then free the packet
			si := &SinkInfo{Packet: rp}
			for _, sk := range s.sink {
				sk.HandleData(s, si)
			}
			rp.FreePacket()
		case <-cancelC:
			return
		}
	}
}

func (s *MediaSession) sendPacketLoop(ctx context.Context) {
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "send"})
	gauge.Inc()

	defer func() {
		if r := recover(); r != nil {
			logger.Errorln("sendPacketLoop panic %v", r)
			debug.PrintStack()
		}
	}()
	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop local send", s.GetSessionId())
		s.doneC <- "done"
	}()

	timeStep := codec.GetCodecTimeStep(s.avPayloadCodec)
	ticker := time.NewTicker(time.Duration(timeStep) * time.Millisecond)
	cancelC := ctx.Done()
outLoop:
	for {
		select {
		case <-ticker.C:
			si := &SourceInfo{}
			// pull data from all sources
			for _, source := range s.source {
				source.PullData(s, si)
			}

			// send all packets based on SourceInfo link list
			// for video, a frame can have more than one packet with same timestamp
			for si != nil {
				payload, pts, mark := si.Payload, si.Pts, si.Marker
				if payload != nil {
					if s.rtpSession == nil {
						break outLoop
					}
					packet := s.rtpSession.NewDataPacket(pts)
					packet.SetMarker(mark)
					packet.SetPayload(payload)
					_, _ = s.rtpSession.WriteData(packet)
					packet.FreePacket()
				}
				si = si.Next
			}
		case <-cancelC:
			break outLoop
		}
	}

	ticker.Stop()
}
