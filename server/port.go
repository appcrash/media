package server

import "github.com/appcrash/media/server/prom"

type portPool struct {
	freePortC  chan uint16 // store rtp ports (even number)
	start, end uint16
}

func (p *portPool) init(start uint16, end uint16) {
	// NOTE: rtp use even port
	if start&0x01 != 0 {
		start += 1
	}
	if end&0x01 != 0 {
		end -= 1
	}
	pairs := (end - start) / 2
	p.freePortC = make(chan uint16, pairs)
	go func() {
		for i := start; i < end; i += 2 {
			p.freePortC <- i
		}
	}()
	p.start = start
	p.end = end
}

func (p *portPool) get() (port uint16) {
	// if pool is empty, return 0 as this port wouldn't be used by applications
	select {
	case port = <-p.freePortC:
		prom.UsedPortPair.Inc()
	default:
	}
	return
}

func (p *portPool) put(port uint16) {
	p.freePortC <- port
	prom.UsedPortPair.Dec()
}
