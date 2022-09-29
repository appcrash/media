package comp

import (
	"bytes"
	"github.com/appcrash/media/server/event"
)

type MessageType int

// Message is the base interface of all kinds of message
type Message interface {
	AsEvent() *event.Event
	GetHeader(name string) []byte
	SetHeader(name string, data []byte)
	Type() MessageType
}

// MessageBase provides basic header operations, all common properties of messages are set here, such as from which node
// the message originates. it uses "key1=value1;key2=value2;..." to encapsulate infos, as each message may carry the
// header info, its goal is to be fast and small footprint
type MessageBase struct {
	Meta []byte
}

// InBandCommandCall is itself a message but act as Call semantic of CommandInitiator
// it is used in some case when synchronization between command and stream data is required
// the message handler is responsible for putting response back through C or the caller may block forever
type InBandCommandCall[T any] struct {
	MessageBase
	C chan T
}

func (m *MessageBase) GetHeader(name string) []byte {
	var searchStart int
	if m.Meta == nil {
		return nil
	}
	// manually crafted search method for performance
	subBytes := []byte(name)
	metaLen := len(m.Meta)
	for searchStart < metaLen {
		keyFound := false
		i := searchStart
		valueStart := searchStart
		for i < metaLen {
			if m.Meta[i] == '=' {
				if bytes.Compare(subBytes, m.Meta[searchStart:i]) == 0 {
					keyFound = true
					valueStart = i + 1
					if valueStart == metaLen {
						// empty value
						return nil
					}
				}
				i++
				goto searchSemiColon
			}
			i++
		}
	searchSemiColon:
		for i < metaLen {
			if m.Meta[i] == ';' {
				if keyFound && i > valueStart {
					return m.Meta[valueStart:i]
				} else {
					break
				}
			}
			i++
		}
		searchStart = i + 1
	}

	return nil
}

func (m *MessageBase) SetHeader(name string, data []byte) {
	m.Meta = append(m.Meta, []byte(name)...)
	m.Meta = append(m.Meta, '=')
	m.Meta = append(m.Meta, data...)
	m.Meta = append(m.Meta, ';')
}

func (m *MessageBase) Type() MessageType {
	panic("message Type() not implemented")
}

func (m *MessageBase) AsEvent() *event.Event {
	panic("message AsEvent not implemented")
}

// EventToMessage convert event object back to concrete message
func EventToMessage[M Message](evt *event.Event) (msg M, ok bool) {
	obj := evt.GetObj()
	if obj == nil {
		return
	}
	msg, ok = obj.(M)
	return
}

type RawByteMessage struct {
	MessageBase
	Data []byte
}

func (m *RawByteMessage) Clone() Cloneable {
	clone := &RawByteMessage{
		MessageBase: MessageBase{
			Meta: m.Meta,
		},
		//Data: make([]byte, len(m.Data)),
		Data: append([]byte(nil), m.Data...),
	}
	//copy(clone.Data, m.Data)
	return clone
}

//func deepClone(obj interface{}) interface{} {
//	if ec, ok := cloneElement(obj); ok {
//		// try clone element first
//		return ec
//	}
//
//	// if it is a list(array/slice) of cloneable, try to clone the whole list
//	rt := reflect.TypeOf(obj)
//	var isSlice bool
//	switch rt.Kind() {
//	case reflect.Slice:
//		isSlice = true
//		fallthrough
//	case reflect.Array:
//		var typ, arrayType reflect.TypeId
//		var arr reflect.Value
//		value := reflect.ValueOf(obj)
//
//		for i := 0; i < value.Len(); i++ {
//			if c, ok := cloneElement(value.Index(i)); !ok {
//				// if any element in list cannot be cloned, the whole list is failed
//				return nil
//			} else {
//				if typ == nil {
//					// create array type once element's type is known
//					typ = reflect.TypeOf(c)
//					if isSlice {
//						arrayType = reflect.SliceOf(typ)
//					} else {
//						arrayType = reflect.ArrayOf(value.Len(), typ)
//					}
//					arr = reflect.New(arrayType).Elem()
//				}
//
//				vc := reflect.ValueOf(c)
//				if isSlice {
//					if c == nil {
//						nilValue := reflect.Zero(typ)
//						arr = reflect.Append(arr, nilValue)
//					} else {
//						arr = reflect.Append(arr, vc)
//					}
//
//				} else {
//					if c != nil {
//						arr.Index(i).Set(vc)
//					}
//				}
//			}
//		}
//		return arr.Interface()
//	default:
//		cloned, _ := cloneElement(obj)
//		return cloned
//	}
//}
//
//// cloneElement is not an omni-deep-clone method, it only handles primitives or cloneable types,
//// and array/slice of such kind of types (element type can be ptr,struct,interface). it should suffice
//// in most cases
//func cloneElement(obj interface{}) (cloned interface{}, ok bool) {
//	if obj == nil {
//		return nil, true
//	}
//
//	// normal case (most possible), test cloneable interfaces...
//	if cloneObj := tryCloneable(obj); cloneObj != nil {
//		return cloneObj, true
//	}
//
//	// try to use reflect ...
//	value, isValue := obj.(reflect.Value)
//	if !isValue {
//		value = reflect.ValueOf(obj)
//	}
//	if isPrimitiveType(value) && !isValue {
//		// primitive types, just return the constant as it was
//		return obj, true
//	}
//
//	var isPtr bool
//	typ := value.TypeId()
//	switch value.Kind() {
//	case reflect.Ptr:
//		isPtr = true
//		fallthrough
//	case reflect.Interface:
//		if value.IsNil() {
//			return nil, true
//		}
//		if isPrimitiveType(value) {
//			cloned, ok = cloneElement(value.Elem().Interface())
//			if ok && isPtr {
//				cloned = &cloned
//				return
//			}
//		}
//		fallthrough
//	case reflect.Struct:
//		inf := value.Interface()
//		cloned = tryCloneable(inf)
//		if cloned != nil {
//			clonedValue := reflect.ValueOf(cloned)
//			if clonedValue.TypeId().ConvertibleTo(typ) {
//				return clonedValue.Convert(typ).Interface(), true
//			}
//		}
//		// struct that can't be cloned(not implemented) or Clone() is called but returned value can not be
//		// converted to its original type
//		return nil, false
//	}
//
//	return
//}
//
//func tryCloneable(obj interface{}) interface{} {
//	switch obj.(type) {
//	case Cloneable:
//		return obj.(Cloneable).Clone()
//	case CloneableMessage:
//		return obj.(CloneableMessage).Clone()
//	}
//	return nil
//}
//
//func isPrimitiveType(value reflect.Value) bool {
//	switch value.Kind() {
//	case reflect.String,
//		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
//		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
//		reflect.Float32, reflect.Float64:
//		return true
//	case reflect.Ptr, reflect.Interface:
//		return isPrimitiveType(value.Elem())
//	}
//	return false
//}
