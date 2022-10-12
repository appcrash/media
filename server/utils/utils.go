package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

func RemoveElementFromArray[T any](slice []T, index int) ([]T, error) {
	if index > len(slice)-1 || index < 0 {
		return nil, fmt.Errorf("index out of range")
	}
	return append(slice[:index], slice[index+1:]...), nil
}

const regSnakeToCamelCasePattern = `_+[a-z]`

var regSnakeToCamelCase = regexp.MustCompile(regSnakeToCamelCasePattern)

// SnakeToCamelCase converts string in form of foo_bar or foo__bar ... into fooBar
func SnakeToCamelCase(str string) string {
	return regSnakeToCamelCase.ReplaceAllStringFunc(str, func(match string) string {
		last := string(match[len(match)-1])
		return strings.ToUpper(last)
	})
}

// CamelCaseToSnake converts string in form of fooBar into foo_bar
func CamelCaseToSnake(str string) string {
	var sb strings.Builder
	var hasSnake = true // avoid inserting '_' at head
	for _, c := range str {
		if c == '_' {
			if !hasSnake {
				sb.WriteByte('_')
			}
			hasSnake = true
		} else {
			if unicode.IsUpper(c) {
				if !hasSnake {
					sb.WriteByte('_')
				}
				sb.WriteRune(unicode.ToLower(c))
			} else {
				sb.WriteRune(c)
			}
			hasSnake = false
		}
	}
	return sb.String()
}

type flagSize interface {
	uint8 | uint16 | uint32 | uint64
}

type Flag[T flagSize] struct {
	Flag T
}

func (f *Flag[T]) SetFlag(bit T) {
	f.Flag |= bit
}

func (f *Flag[T]) ClearFlag(bit T) {
	f.Flag &= ^bit
}

func (f *Flag[T]) HasFlag(bit T) bool {
	return f.Flag&bit != 0
}

// WaitChannelWithTimeout read exactly number of nbWaitFor values out from channel then return
// if this cannot achieve within timeout duration, non-nil error returned
func WaitChannelWithTimeout[T any](c <-chan T, nbWaitFor int, d time.Duration) (err error) {
	if nbWaitFor <= 0 {
		return
	}
	ready := 0
	timeoutC := time.After(d)
	for ready < nbWaitFor {
		select {
		case <-c:
			ready++
		case <-timeoutC:
			err = errors.New("waiting channel times out")
			return
		}
	}
	return
}

// AopCall take a pointer to struct (obj) and check all embedded(of type struct) VISIBLE FIELD'S pointer type,
// if this type implements the provided interface type, invoke one of its method provided by the method name, finally if
// the struct itself implements the interface, invoke for it too. this is different to golang's method shadowing spec,
// which says only the top-most method(the one shadowing same method of its parent) is seen except explicitly invoking
// parent's method. this function emulates aop functionality, when an interface method implemented by an object is
// invoked(point cut, before or after), execute other operations just by embedding corresponding atomic struct.
// currently AopCall only supports before pointcut
func AopCall(obj interface{}, args []interface{}, interfaceType reflect.Type, methodName string) (rv [][]reflect.Value) {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type().Elem()

	if objType.Kind() != reflect.Struct {
		// obj is not a ptr to a struct
		return
	}
	if interfaceType.Kind() != reflect.Interface {
		return
	}
	methodType, exist := interfaceType.MethodByName(methodName)
	if !exist {
		return
	}
	var ptrType reflect.Type
	var ptrValue reflect.Value
	// check every embedded field of a struct and invoke interface method for them
	for i := 0; i <= objType.NumField(); i++ {
		if i == objType.NumField() {
			// check the object itself
			ptrValue = objValue
			ptrType = objValue.Type()
		} else {
			// check the embedded field
			ptrValue = objValue.Elem().Field(i).Addr()
			field := objType.Field(i)
			fieldType := field.Type
			if !field.IsExported() || fieldType.Kind() != reflect.Struct {
				continue
			}
			ptrType = reflect.PtrTo(fieldType)
		}

		if ptrType.Implements(interfaceType) {
			method := ptrValue.Convert(interfaceType).Method(methodType.Index)
			var valueArgs []reflect.Value
			for _, a := range args {
				valueArgs = append(valueArgs, reflect.ValueOf(a))
			}
			rv = append(rv, method.Call(valueArgs))
		}
	}
	return
}

// SetField set a struct field even it is not exported
func SetField(field, value reflect.Value) {
	if field.CanSet() {
		field.Set(value)
	} else {
		// forcefully setting unexported variable
		nf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		nf.Set(value)
	}
}
