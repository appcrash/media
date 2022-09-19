package comp

type graphCommandInitiator struct {
	composer *Composer
}

func (g *graphCommandInitiator) Call(fromNode, toNode string, args []string) (resp []string) {
	if fromNode == toNode {
		return WithError("can not call to self")
	}
	to := g.composer.GetNode(toNode)
	if to == nil {
		return WithError("to node not exist")
	}
	resp = to.OnCall(fromNode, args)

	return
}

func (g *graphCommandInitiator) Cast(fromNode, toNode string, args []string) {
	if fromNode == toNode {
		return
	}
	to := g.composer.GetNode(toNode)
	if to == nil {
		return
	}
	to.OnCast(fromNode, args)
}

// InitiatorNode is used when node requires to send commands
type InitiatorNode struct {
	nodeName  string
	initiator CommandInitiator
}

func (i *InitiatorNode) BeforeCompose(c *Composer, node SessionAware) error {
	i.initiator = c.GetCommandInitiator()
	i.nodeName = node.GetNodeName()
	return nil
}

func (i *InitiatorNode) Call(toNode string, args []string) (resp []string) {
	return i.initiator.Call(i.nodeName, toNode, args)
}

func (i *InitiatorNode) Cast(toNode string, args []string) {
	i.initiator.Cast(i.nodeName, toNode, args)
}
