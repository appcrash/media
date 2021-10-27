package channel_test

import (
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/rpc"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

const testInstanceId = "testInstance"

func init() {
	channel.InitLogger(logrus.New())
}

func TestKeepaliveTimeout(t *testing.T) {
	sc := channel.GetSystemChannel()
	instanceState, err := sc.RegisterInstance(testInstanceId)
	if err != nil {
		t.Fatal("register instance failed")
	}
	select {
	case <-instanceState.FromInstanceC:
	case <-time.After(channel.KeepAliveTimeout + 1*time.Second):
		t.Fatal("the instance should time out")
	}
	msg := rpc.SystemEvent{
		Cmd:        rpc.SystemCommand_SESSION_INFO,
		InstanceId: testInstanceId,
		SessionId:  "some_session",
		Event:      "start",
	}
	if sc.NotifyInstance(&msg) == nil {
		t.Fatal("notify disconnected instance should return error")
	}
}

func ExampleBroadcast() {
	const n = 5
	wg := &sync.WaitGroup{}
	wg.Add(n)
	sc := channel.GetSystemChannel()
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("%v_%d", testInstanceId, i)
		is, _ := sc.RegisterInstance(id)
		go func() {
			select {
			case msg := <-is.ToInstanceC:
				fmt.Printf("msg from %v\n", msg.InstanceId)
			}
			wg.Done()
		}()
	}
	for i := 0; i < n; i++ {
		msg := rpc.SystemEvent{
			Cmd:        rpc.SystemCommand_SESSION_INFO,
			InstanceId: fmt.Sprintf("%v_%d", testInstanceId, i),
			SessionId:  "some_session_id",
			Event:      "start ",
		}
		sc.NotifyInstance(&msg)
	}
	wg.Wait()

	// Unordered OUTPUT:
	// msg from testInstance_0
	// msg from testInstance_1
	// msg from testInstance_2
	// msg from testInstance_3
	// msg from testInstance_4
}
