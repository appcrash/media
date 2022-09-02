package server

import (
	"sync"
)

func newPortPool() *portPool {
	return &portPool{
		freePortSet: make(map[uint16]struct{}),
	}
}

type portPool struct {
	mutex sync.Mutex

	freePortSet map[uint16]struct{} // store rtp ports (even number)
	start, end  uint16
}

func (p *portPool) init(start uint16, end uint16) {
	// NOTE: rtp use even port
	if start&0x01 != 0 {
		start += 1
	}
	if end&0x01 != 0 {
		end -= 1
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for i := start; i < end; i += 2 {
		p.freePortSet[i] = struct{}{}
	}

	p.start = start
	p.end = end
}

func (p *portPool) get() (port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// if pool is empty, return 0 as this port wouldn't be used by applications

	if len(p.freePortSet) == 0 {
		goto noPort
	}
	for port, _ = range p.freePortSet {
		break
	}
	delete(p.freePortSet, port)

noPort:
	return
}

func (p *portPool) put(port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.freePortSet[port]; ok {
		logger.Errorf("put port: %v to pool but it is already in it", port)
	}
	p.freePortSet[port] = struct{}{}
}
