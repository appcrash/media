package comp

type RtpSrc struct {
	SessionNode
}

func (s *RtpSrc) handleA(m *RawByteMessage) {

}

func (s *RtpSrc) Accept() []MessageType {
	return nil
}
