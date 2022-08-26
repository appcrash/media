package comp

import "reflect"

// MessageTrait is used to record all possible message kinds that flow among nodes of known types, ensure node links
// are compatible, that is Node A output is accepted by Node B input if there would be a link.
type MessageTrait struct {
	Name string
	Type reflect.Type
}

var messageTraitRegistry map[string]*MessageTrait

func RegisterMessageType(model Message) {
	typ := reflect.TypeOf(model)
	name := model.Name()
	mt := &MessageTrait{
		Name: name,
		Type: typ,
	}
	logger.Infof("register message type (Name:%v) => (type:%v)", name, typ.Name())
	messageTraitRegistry[name] = mt
}

func MessageTraitOf(model Message) (exist bool, mt MessageTrait) {
	var pmt *MessageTrait
	name := model.Name()
	if pmt, exist = messageTraitRegistry[name]; exist {
		mt = *pmt
	}
	return
}
