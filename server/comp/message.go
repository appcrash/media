package comp

import (
	"github.com/appcrash/media/server/event"
)

type MessageType int

// Message is the base interface of all kinds of message
type Message interface {
	AsEvent() *event.Event
	GetMeta() []byte
	Type() MessageType
}

type Cloneable interface {
	Clone() Cloneable
}

type BaseMessage struct {
	TypeId MessageType
	Meta   []byte
}

func (m *BaseMessage) AsEvent() *event.Event {
	return event.NewEvent(int(m.Type()), m)
}

func (m *BaseMessage) GetMeta() []byte {
	return m.Meta
}

func (m *BaseMessage) Type() MessageType {
	return m.TypeId
}

type MatchAnyMessage struct {
	BaseMessage
}

func (m *MatchAnyMessage) Type() MessageType {
	return ANY
}

type RawByteMessage struct {
	BaseMessage
	Data []byte
}

func (m *RawByteMessage) Type() MessageType {
	return RawByte
}

func (m *RawByteMessage) Clone() Cloneable {
	clone := &RawByteMessage{
		BaseMessage: BaseMessage{
			TypeId: m.TypeId,
			Meta:   m.Meta,
		},
		Data: make([]byte, len(m.Data)),
	}
	copy(clone.Data, m.Data)
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
