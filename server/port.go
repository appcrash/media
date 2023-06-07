package server

import (
	"github.com/appcrash/media/server/utils"
	"sync"
)

func NewPortPool() *PortPool {
	return &PortPool{
		freePortSet: utils.NewSet[uint16](),
	}
}

type PortPool struct {
	mutex sync.Mutex

	freePortSet *utils.Set[uint16] // store rtp ports (even number)
	start, end  uint16
}

func (p *PortPool) Init(start uint16, end uint16) {
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
		p.freePortSet.Add(i)
	}

	p.start = start
	p.end = end
}

func (p *PortPool) Get() (port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// if pool is empty, return 0 as this port wouldn't be used by applications

	if p.freePortSet.Size() == 0 {
		goto noPort
	}
	port = p.freePortSet.GetAndRemove()
noPort:
	return
}

func (p *PortPool) Put(port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.freePortSet.Contain(port) {
		logger.Errorf("put port: %v to pool but it is already in it", port)
		return
	}
	p.freePortSet.Add(port)
}
