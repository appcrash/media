package comp

// ChannelSink forward input message to go channel
type ChannelSink struct {
	SessionNode

	channel string // channel name to be linked
}

func NewChannelSink() SessionAware {
	c := &ChannelSink{}
	return c
}
