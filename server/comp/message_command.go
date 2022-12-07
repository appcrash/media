package comp

// all messages below is effectively commands for nodes

// LinkPointRequestMessage received when another node wants connecting to me
type LinkPointRequestMessage struct {
	InBandCommandCall[*MessageTrait]
	PreferredTrait []*MessageTrait
	LinkIdentity   LinkIdentityType
}

// ChannelLinkRequestMessage received when being ask to link to a provided channel
type ChannelLinkRequestMessage struct {
	InBandCommandCall[interface{}]
	LinkChannel chan []byte
}
