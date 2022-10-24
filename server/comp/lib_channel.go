package comp

import (
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/rpc"
)

// ChannelNode enables a node to send async event to signalling service through system channel
type ChannelNode struct {
	sessionId, instanceId, nodeName string
}

func (n *ChannelNode) BeforeCompose(c *Composer, node SessionAware) error {
	n.sessionId = c.GetSessionId()
	n.instanceId = c.GetInstanceId()
	n.nodeName = node.GetNodeName()
	return nil
}

func (n *ChannelNode) NotifyInstance(event string) error {
	return channel.GetSystemChannel().NotifyInstance(&rpc.SystemEvent{
		Cmd:        rpc.SystemCommand_USER_EVENT,
		InstanceId: n.instanceId,
		SessionId:  n.sessionId,
		Event:      event,
	})
}

func (n *ChannelNode) BroadcastInstance(event string) error {
	return channel.GetSystemChannel().BroadcastInstance(&rpc.SystemEvent{
		Cmd:        rpc.SystemCommand_USER_EVENT,
		InstanceId: n.instanceId,
		SessionId:  n.sessionId,
		Event:      event,
	})
}
