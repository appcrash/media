package comp

// predefined message types

const (
	MtNewLinkPoint = iota
	MtConnectNode
	MtChannelLink
	MtRawByte
	MtUserMessageStart
)

type LinkPointMessage struct {
	InBandCommandCall[*MessageTrait]
	OfferedTrait []*MessageTrait
	LinkIdentity LinkIdentityType
}

func (l *LinkPointMessage) Type() MessageType {
	return MtNewLinkPoint
}

type ConnectNodeMessage struct {
	InBandCommandCall[bool]
	Session, NodeName string
}

func (c *ConnectNodeMessage) Type() MessageType {
	return MtConnectNode
}

type ChannelLinkMessage struct {
	InBandCommandCall[interface{}]
	LinkChannel chan []byte
}

func (c *ChannelLinkMessage) Type() MessageType {
	return MtChannelLink
}
