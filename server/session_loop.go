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
					s.Stop()
					return
				}
			}
		case <-ctx.Done():
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
	var data []byte
	for {
		select {
		case rp, more := <-dataReceiver:
			var shouldContinue bool
			if !more {
				// RTP stack closed this channel, so stop receiving anymore
				return
			}

			// push received data to all sinks, then free the packet
			data = nil
			for _, sk := range s.sink {
				data, shouldContinue = sk.HandleData(s, rp, data)
				if !shouldContinue {
					break
				}
			}
			rp.FreePacket()
		case <-ctx.Done():
			return
		}
	}
}

func (s *MediaSession) sendPacketLoop(ctx context.Context) {
	var ts uint32 = 0
	gauge := prom.SessionGoroutine.With(prometheus.Labels{"type": "send"})
	gauge.Inc()

	defer func() {
		if r := recover(); r != nil {
			logger.Fatalln("sendPacketLoop panic %v", r)
			debug.PrintStack()
		}
	}()
	defer func() {
		gauge.Dec()
		logger.Debugf("session:%v stop local send", s.GetSessionId())
		s.doneC <- "done"
	}()

	timeStep := codec.GetCodecTimeStep(s.audioPayloadCodec)
	ticker := time.NewTicker(time.Duration(timeStep) * time.Millisecond)

outLoop:
	for {
		select {
		case <-ticker.C:
			var data []byte
			var tsDelta uint32

			// pull data from all sources
			for _, source := range s.source {
				data, tsDelta = source.PullData(s, data, tsDelta)
			}
			if data != nil {
				if s.rtpSession == nil {
					break outLoop
				}
				packet := s.rtpSession.NewDataPacket(ts)
				packet.SetPayload(data)
				_, _ = s.rtpSession.WriteData(packet)
				packet.FreePacket()
				ts += tsDelta
			}
		case <-ctx.Done():
			break outLoop
		}
	}

	ticker.Stop()
}
