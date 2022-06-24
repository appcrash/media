package comp

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"reflect"
)

// RawByteMessage is used to pass byte streaming data between nodes intra/inter session, for efficiency
type RawByteMessage []byte

// CtrlMessage is used to invoke or cast function call
type CtrlMessage struct {
	M []string
	C chan []string // used to receive result if not nil
}

// GenericMessage is used to encapsulate custom object in message with tagged type name
// nodes can only communicate with each other who can build/handle GenericMessage of the same tagged type
type GenericMessage struct {
	Subtype string
	Obj     interface{}
}

// Message is the base interface of all kinds of message
type Message interface {
	AsEvent() *event.Event
}

type Cloneable interface {
	Clone() Cloneable
}

// CloneableMessage is a must if message need to pass through pubsub node
type CloneableMessage interface {
	Message
	Clone() CloneableMessage
}

func NewRawByteMessage(d string) RawByteMessage {
	return RawByteMessage(d)
}

func (m RawByteMessage) String() string {
	return string(m)
}
func (m RawByteMessage) Clone() CloneableMessage {
	mc := make(RawByteMessage, len(m))
	copy(mc, m)
	return mc
}

func (m RawByteMessage) AsEvent() *event.Event {
	return event.NewEvent(RawByte, m)
}

func (cm *CtrlMessage) AsEvent() *event.Event {
	if cm.C != nil {
		return event.NewEvent(CtrlCall, cm)
	} else {
		return event.NewEvent(CtrlCast, cm)
	}
}

func (gm *GenericMessage) AsEvent() *event.Event {
	return event.NewEvent(Generic, gm)
}

func (gm *GenericMessage) String() string {
	return fmt.Sprintf("GenericMessage type:%v value:%v", gm.Subtype, gm.Obj)
}

// Clone returns non-nil object only if internal object is also cloneable or an array(slice) of cloneable
func (gm *GenericMessage) Clone() (obj CloneableMessage) {
	if gm.Obj == nil || gm.Obj == gm {
		// prevent recursive clone
		return
	}

	cloned := deepClone(gm.Obj)
	if cloned == nil {
		return nil
	}
	obj = &GenericMessage{
		Subtype: gm.Subtype,
		Obj:     cloned,
	}
	return
}

// one-level deep clone
func deepClone(obj interface{}) interface{} {
	if ec := cloneElement(obj); ec != nil {
		// try clone element first
		return ec
	}

	// if it is a list(array/slice) of cloneable, try to clone the whole list
	rt := reflect.TypeOf(obj)
	var isSlice bool
	switch rt.Kind() {
	case reflect.Slice:
		isSlice = true
		fallthrough
	case reflect.Array:
		var typ, arrayType reflect.Type
		var arr reflect.Value
		value := reflect.ValueOf(obj)

		for i := 0; i < value.Len(); i++ {
			if c := cloneElement(value.Index(i)); c == nil {
				// if any element in list cannot be cloned, the whole list is failed
				return nil
			} else {
				if typ == nil {
					// create array type once element's type is known
					typ = reflect.TypeOf(c)
					if isSlice {
						arrayType = reflect.SliceOf(typ)
					} else {
						arrayType = reflect.ArrayOf(value.Len(), typ)
					}
					arr = reflect.New(arrayType).Elem()
				}

				if isSlice {
					arr = reflect.Append(arr, reflect.ValueOf(c))
				} else {
					arr.Index(i).Set(reflect.ValueOf(c))
				}
			}
		}
		return arr.Interface()
	default:
		return cloneElement(obj)
	}
}

func cloneElement(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	// normal case (most possible), test cloneable interfaces...
	if cloneObj := tryCloneable(obj); cloneObj != nil {
		return cloneObj
	}

	// try to use reflect ...
	value, isValue := obj.(reflect.Value)
	if !isValue {
		value = reflect.ValueOf(obj)
	}
	if isPrimitiveType(value) && !isValue {
		// primitive types, just return the constant as it was
		return obj
	}

	var isPtr bool
	var clonedObj interface{}
	typ := value.Type()
	switch value.Kind() {
	case reflect.Ptr:
		isPtr = true
		fallthrough
	case reflect.Interface:
		if isPrimitiveType(value) {
			clonedObj = cloneElement(value.Elem().Interface())
			if isPtr {
				clonedObj = &clonedObj
			}
		}
		fallthrough
	case reflect.Struct:
		inf := value.Interface()
		clonedObj = tryCloneable(inf)
		if clonedObj != nil {
			newValue := reflect.ValueOf(clonedObj).Convert(typ)
			return newValue.Interface()
		}

	}

	return clonedObj
}

func tryCloneable(obj interface{}) interface{} {
	// normal case (most possible), test cloneable interfaces...
	switch obj.(type) {
	case Cloneable:
		return obj.(Cloneable).Clone()
	case CloneableMessage:
		return obj.(CloneableMessage).Clone()
	}
	return nil
}

func isPrimitiveType(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Ptr, reflect.Interface:
		return isPrimitiveType(value.Elem())
	}
	return false
}

//cloneableType := reflect.TypeOf((*Cloneable)(nil)).Elem()
