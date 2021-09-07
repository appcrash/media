package comp

// Controller send CtrlMessage to nodes and get replies from them, and simply use string as rpc protocol.
// this file provides utility functions to build and parse strings based on the protocol.

func With(args ...string) []string {
	return args
}

// WithConnect dynamically ask the callee to set data pipe endpoint
func WithConnect(toSession, toName string) []string {
	return With("conn", toSession, toName)
}

// WithOk return normal reply
func WithOk(args ...string) (r []string) {
	r = append(r, "ok")
	r = append(r, args...)
	return
}

// WithError return abnormal reply
func WithError(args ...string) (r []string) {
	r = append(r, "err")
	r = append(r, args...)
	return
}
