package server

import (
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/codec"
	"runtime/debug"
	"time"
)

// receive rtcp packet
func (s *MediaSession) receiveCtrlLoop() {
	rtcpReceiver := s.rtpSession.CreateCtrlEventChan()
	ctrlC := s.rcvRtcpCtrlC

	defer func() {
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
		case msg := <-ctrlC:
			if msg == "stop" {
				return
			}
		}
	}
}

func (s *MediaSession) receivePacketLoop() {
	// Create and store the data receive channel.
	defer func() {
		if r := recover(); r != nil {
			logger.Fatalln("receivePacketLoop panic(recovered)")
			debug.PrintStack()
		}
	}()

	defer func() {
		s.doneC <- "done"
	}()

	rtpSession := s.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()
	var data []byte
outLoop:
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
		case cmd := <-s.rcvCtrlC:
			if cmd == "stop" {
				logger.Debugf("session:%v stop local receive", s.GetSessionId())
				break outLoop
			}
		}
	}

}

func (s *MediaSession) sendPacketLoop() {
	var ts uint32 = 0
	defer func() {
		if r := recover(); r != nil {
			logger.Fatalln("sendPacketLoop panic %v", r)
			debug.PrintStack()
		}
	}()
	defer func() {
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
		case cmd := <-s.sndCtrlC:
			if cmd == "stop" {
				logger.Debugf("session:%v stop local send", s.GetSessionId())
				break outLoop
			}
		}
	}

	ticker.Stop()
}
