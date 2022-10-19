package channel

import (
	"fmt"
	"github.com/appcrash/media/server/rpc"
	"sync"
	"time"
)

type Listener interface {
	// OnChannelEvent invoked by multiple goroutine concurrently
	OnChannelEvent(e *rpc.SystemEvent)
}

type InstanceState struct {
	name                       string
	lastSeen                   time.Time
	FromInstanceC, ToInstanceC chan *rpc.SystemEvent
}

// Channel is a bidirectional channel that communicates with each signalling server (by grpc streaming)
// the goal is to provide mechanism enabling media server to send async event to other services and
// route incoming async events to subsystem in media server
type Channel struct {
	mutex            sync.Mutex
	instanceStateMap map[string]*InstanceState
	listeners        []Listener
}

// singleton sys channel
var sysChannel = newChannel()

const (
	KeepAliveCheckDuration = 2 * time.Second
	KeepAliveTimeout       = KeepAliveCheckDuration * 3
)

func newChannel() *Channel {
	return &Channel{
		instanceStateMap: make(map[string]*InstanceState),
	}
}

func GetSystemChannel() *Channel {
	return sysChannel
}

// close must be called with mutex held
func (is *InstanceState) close() {
	from := is.FromInstanceC
	to := is.ToInstanceC
	is.FromInstanceC = nil
	if from != nil {
		close(from) // end channel receive loop
	}
	is.ToInstanceC = nil
	if to != nil {
		close(to) // end rpc receive loop
	}
}

func (sc *Channel) AddListener(l Listener) {
	sc.mutex.Lock()
	sc.listeners = append(sc.listeners, l)
	sc.mutex.Unlock()
}

func (sc *Channel) RegisterInstance(name string) (is *InstanceState, err error) {
	if name == "" {
		err = fmt.Errorf("invalid instance name")
		return
	}

	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if prevIs, exist := sc.instanceStateMap[name]; exist {
		// an instance(with same id) disconnected and reconnects to media server
		// previous instance state's channels can be used concurrently
		// even though it is deleted from map right now. set send-channel to nil to stop send, and close recv channel.
		// set to nil before closing
		prevIs.close()
		delete(sc.instanceStateMap, name)
	}

	is = &InstanceState{
		name:          name,
		lastSeen:      time.Now(),
		FromInstanceC: make(chan *rpc.SystemEvent, 64),
		ToInstanceC:   make(chan *rpc.SystemEvent, 64),
	}
	sc.instanceStateMap[name] = is
	go sc.startReceiveLoop(name)
	return
}

func (sc *Channel) HasInstance(name string) (exist bool) {
	sc.mutex.Lock()
	_, exist = sc.instanceStateMap[name]
	sc.mutex.Unlock()
	return
}

// NotifyInstance NONBLOCK send event to instance
func (sc *Channel) NotifyInstance(se *rpc.SystemEvent) (err error) {
	if se.InstanceId == "" {
		return fmt.Errorf("invalid instance id when notifying instance")
	}
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if is, exist := sc.instanceStateMap[se.InstanceId]; exist {
		select {
		case is.ToInstanceC <- se:
		default:
			err = fmt.Errorf("server channel: send to instance %v failed", se.InstanceId)
		}
	} else {
		err = fmt.Errorf("server channel: no such instance %v when send to instance", se.InstanceId)
	}
	return
}

// BroadcastInstance NONBLOCK send event to all instances
func (sc *Channel) BroadcastInstance(se *rpc.SystemEvent) (err error) {
	if se.InstanceId == "" {
		return fmt.Errorf("invalid instance id when broadcasting instances")
	}
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	for name, is := range sc.instanceStateMap {
		select {
		case is.ToInstanceC <- se:
		default:
			err = fmt.Errorf("server channel: broadcast to instance %v failed", name)
			return
		}
	}
	return
}

func (sc *Channel) startReceiveLoop(instanceId string) {
	var is *InstanceState
	var exist bool
	sc.mutex.Lock()
	if is, exist = sc.instanceStateMap[instanceId]; !exist {
		logger.Errorf("server channel recv loop with invalid instance id:%v", instanceId)
		sc.mutex.Unlock()
		return
	}
	sc.mutex.Unlock()
	ticker := time.NewTicker(KeepAliveCheckDuration)
	for {
		select {
		case se, more := <-is.FromInstanceC:
			if !more {
				logger.Infoln("from-instance channel closed")
				return
			}
			handled := true
			switch se.Cmd {
			case rpc.SystemCommand_KEEPALIVE:
				is.lastSeen = time.Now()
				rse := rpc.SystemEvent{
					Cmd:        rpc.SystemCommand_KEEPALIVE,
					InstanceId: se.InstanceId,
					SessionId:  se.SessionId,
					Event:      se.Event,
				}
				sc.NotifyInstance(&rse)
			default:
				handled = false
			}

			if handled {
				continue
			}
			is.lastSeen = time.Now()

			// invoke listeners
			sc.mutex.Lock()
			listeners := sc.listeners
			sc.mutex.Unlock()
			for _, l := range listeners {
				l.OnChannelEvent(se)
			}
		case <-ticker.C:
			if time.Since(is.lastSeen) > KeepAliveTimeout {
				logger.Errorf("server channel for instance(%v) keep-alive times out, exit recv loop", instanceId)
				sc.mutex.Lock()
				is.close()
				sc.mutex.Unlock()
				return
			}
		}

	}
}
