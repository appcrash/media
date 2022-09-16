package comp

// all messages below is effectively commands for nodes

// LinkPointRequestMessage received when another node wants connecting to me
type LinkPointRequestMessage struct {
	InBandCommandCall[*MessageTrait]
	OfferedTrait []*MessageTrait
	LinkIdentity LinkIdentityType
}

// ConnectNodeRequestMessage received when being ask to connect to another node
type ConnectNodeRequestMessage struct {
	InBandCommandCall[bool]
	Session, NodeName    string
	PreferredMessageName []string // offer
}

// ChannelLinkRequestMessage received when being ask to link to a provided channel
type ChannelLinkRequestMessage struct {
	InBandCommandCall[interface{}]
	LinkChannel chan []byte
}
