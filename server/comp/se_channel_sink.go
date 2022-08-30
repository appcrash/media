package comp

// ChanSink forward input message to go channel
type ChanSink struct {
	SessionNode

	channel string // channel name to be linked
}

func NewChannelSink() SessionAware {
	c := &ChanSink{}
	return c
}
