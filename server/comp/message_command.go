package comp

type LinkPointMessage struct {
	InBandCommandCall[*MessageTrait]
	OfferedTrait []*MessageTrait
	LinkIdentity LinkIdentityType
}

type ConnectNodeMessage struct {
	InBandCommandCall[bool]
	Session, NodeName string
}

type ChannelLinkMessage struct {
	InBandCommandCall[interface{}]
	LinkChannel chan []byte
}

func (m *ChannelLinkMessage) AsRawByteMessage() (r *RawByteMessage) {
	return
}
