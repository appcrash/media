package comp

import "strconv"

// builtinCommandInitiator provides default common CALL commands for all nodes:
// -----------------------------------------------------------------------
// conn {local_node_name}  # return link_index if succeed
// conn {session} {node_name} # return link_link if succeed
// enable_link {link_index}
// disable_link {link_index}
type builtinCommandInitiator struct {
	composer *Composer
}

func (c *builtinCommandInitiator) Call(fromNode, toNode string, args []string) (resp []string) {
	if fromNode == toNode {
		return WithError("can not call to self")
	}
	to := c.composer.GetNode(toNode)
	if to == nil {
		return WithError("to node not exist")
	}
	if resp = c.checkBuiltinCall(to, args); resp == nil {
		resp = to.OnCall(fromNode, args)
	}

	return
}

func (c *builtinCommandInitiator) Cast(fromNode, toNode string, args []string) {
	if fromNode == toNode {
		return
	}
	to := c.composer.GetNode(toNode)
	if to == nil {
		return
	}
	to.OnCast(fromNode, args)
}

func (c *builtinCommandInitiator) checkBuiltinCall(to SessionAware, args []string) []string {
	var isEnable bool
	argLen := len(args)
	if argLen == 0 {
		return nil
	}
	switch args[0] {
	case "conn":
		var session, name string
		if argLen == 2 {
			name = args[1]
		} else if argLen == 3 {
			session, name = args[1], args[2]
		} else {
			logger.Errorf("wrong conn command: %v", args)
			return WithError("wrong conn command")
		}
		if len(session) == 0 {
			session = c.composer.sessionId
		}

		if lp, err := to.StreamTo(session, name, to.Offer()); err == nil {
			return WithOk(strconv.Itoa(lp.LinkId()))
		} else {
			return WithError(err.Error())
		}
	case "enable_link":
		isEnable = true
		fallthrough
	case "disable_link":
		if argLen != 2 {
			return WithError("wrong enable/disable link command")
		}
		linkId, err := strconv.Atoi(args[1])
		if err != nil {
			return WithError(err.Error())
		}
		if lp := to.GetLinkPoint(linkId); lp != nil {
			lp.SetEnabled(isEnable)
			return WithOk()
		} else {
			return WithError("wrong link id for enable/disable link command")
		}

	}
	return nil
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
