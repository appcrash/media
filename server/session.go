package server

import (
	"fmt"
	"github.com/go-audio/wav"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"io/ioutil"
	"math"
	"net"
	"os"
	"time"
)

type MediaSession struct {
	sessionId  string
	rtpSession *rtp.Session
	sndCtrlC   chan string
	rcvCtrlC   chan string
}

var sessionMap = cmap.New()

const PCMAPayLoadLength = 160

func WavPayload(wavData []byte) (rtp_payload []byte) {
	start := 0
	for i := 36; i < len(wavData); i++ {
		if string(wavData[i:i+4]) == "data" {
			start = i + 8
			break
		}
	}
	if start == 0 {
		fmt.Errorf("data ERROR")
		return rtp_payload
	}
	rtp_payload = wavData[start:]
	return rtp_payload
}


func createSession(localPort int) *MediaSession{
	tpLocal, _ := rtp.NewTransportUDP(local, localPort,"")
	session := rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, _ := session.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   local.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone : "",
	}, 0, 0)
	session.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(8)

	ms := MediaSession{
			sessionId:  uuid.New().String(),
			rtpSession: session,
			sndCtrlC:   make(chan string,2),  // use buffered version avoiding deadlock
			rcvCtrlC:   make(chan string,2),
	}
	sessionMap.Set(ms.sessionId,&ms)
	return &ms
}


func (session *MediaSession) receivePacketLocal() {
	// Create and store the data receive channel.
	rtpSession := session.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()
	var cnt int
	of,_ := os.Create("record.wav")
	encoder := wav.NewEncoder(of,8000,8,1,6)
	for {
		select {
		case rp := <-dataReceiver:
			if (cnt % 50) == 0 {
				println("Remote receiver got:", cnt, "packets")
			}
			data := rp.Payload()
			for _,d := range data {
				encoder.WriteFrame(d)
			}
			cnt++
			rp.FreePacket()
		case cmd := <-session.rcvCtrlC:
			if cmd == "stop" {
				println("stop local receive")
				encoder.Close()
			}
			return
		}
	}
}

func (session *MediaSession) sendPacketLocal() {
	sendBuffer := []byte{}
	sendPtr := int(math.MaxInt32)
	sendStep := int(0)
	var ts uint32 = 0

	ticker := time.NewTicker(20 * time.Millisecond)

	for {
		select {
			case <- ticker.C:
				// if send buffer is done, stop the ticker until send buffer activated again
				// otherwise forward the send point and feed the data to rtp
				if sendPtr >= len(sendBuffer) {
					fmt.Println("sendPtr out of buffer range, stop ticker")
					ticker.Stop()
				} else {
					nextPtr := sendPtr + sendStep
					if nextPtr > len(sendBuffer) {
						nextPtr = len(sendBuffer)
					}
					payload := sendBuffer[sendPtr : nextPtr]
					packet := session.rtpSession.NewDataPacket(ts)
					ts += uint32(sendStep)  // TODO: use real timestamp step
					packet.SetPayload(payload)
					session.rtpSession.WriteData(packet)
					packet.FreePacket()
					sendPtr = nextPtr
				}
			case cmd := <-session.sndCtrlC:
				if cmd == "stop" {
					break
				} else {
					if mediaData, err := ioutil.ReadFile(cmd); err == nil {
						fmt.Println("start playing file: %v",cmd)
						wavData := WavPayload(mediaData)
						sendBuffer = wavData
						sendPtr = 0
						sendStep = PCMAPayLoadLength
						ticker.Reset(20 * time.Millisecond)
					} else {
						fmt.Println("error: can not find file %v when change send buffer",cmd)
					}
				}
		}
	}

}

func (session *MediaSession) StartSession() {
	session.rtpSession.StartSession()
	go session.receivePacketLocal()
	go session.sendPacketLocal()
}

func (session *MediaSession) AddRemote(ip string, port int) {
	ipaddr := net.ParseIP(ip)

	println("peer ip port ",ipaddr,port)
	println("session is ", session)
	session.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   ipaddr,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	})
}

