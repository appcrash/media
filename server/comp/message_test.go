package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	//logger.SetLevel(logrus.DebugLevel)
	comp.InitLogger(logger)
}

type testCloneableObjPtr struct {
	a int
}

type testCloneableObjStruct struct {
	a int
}

type testNonCloneableObj struct{}

func (t *testCloneableObjPtr) Clone() comp.Cloneable {
	return &testCloneableObjPtr{t.a}
}

func (t testCloneableObjStruct) Clone() comp.Cloneable {
	return testCloneableObjStruct{t.a}
}

func TestGenericMessage_Cloneable(t *testing.T) {
	rbm := comp.RawByteMessage("some_byte")
	gm := &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     rbm,
	}
	cgm := gm.Clone()
	if cgm == nil {
		t.Fatal("generic message does not clone its internal object")
	}
	gm.Obj = gm
	cgm = gm.Clone()
	if cgm != nil {
		t.Fatal("generic message allow recursive clone")
	}
	gm.Obj = nil
	cgm = gm.Clone()
	if cgm != nil {
		t.Fatal("generic message cloned when object is nil")
	}
}

func TestGenericMessage_NonCloneable(t *testing.T) {
	gm := &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     &testNonCloneableObj{},
	}
	cgm := gm.Clone()
	if cgm != nil {
		t.Fatal("should not clone it")
	}
	gm = &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     testNonCloneableObj{},
	}
	cgm = gm.Clone()
	if cgm != nil {
		t.Fatal("should not clone it")
	}
}

func TestGenericMessage_ArraySlice(t *testing.T) {
	testPtrArray := [3]*testCloneableObjPtr{
		&testCloneableObjPtr{1},
		&testCloneableObjPtr{2},
		&testCloneableObjPtr{3},
	}

	testStructArray := [3]testCloneableObjStruct{
		{1},
		{2},
		{3},
	}

	// ############# test slice #############
	gmList := &comp.GenericMessage{
		Subtype: "test_clone_slice_ptr",
		Obj:     testPtrArray[:],
	}
	cgm := gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message with slice failed")
	}
	listObj := cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message with slice failed, object not cloned")
	}
	ctl, ok := listObj.([]*testCloneableObjPtr)
	if !ok {
		t.Fatalf("clone generic message with slice failed, must be slice type, which is %v", reflect.TypeOf(listObj))
	}
	if len(ctl) != 3 {
		t.Fatalf("clone generic message with slice failed, object length not correct, which is %v,type is %v",
			len(ctl), reflect.TypeOf(ctl))
	}

	gmList = &comp.GenericMessage{
		Subtype: "test_clone_slice_struct",
		Obj:     testStructArray[:],
	}
	cgm = gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message(struct) with slice failed")
	}
	listObj = cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message(struct) with slice failed, object not cloned")
	}
	ctlObj, ok1 := listObj.([]testCloneableObjStruct)
	if !ok1 {
		t.Fatal("clone generic message(struct) with slice failed, must be slice type")
	}
	if len(ctlObj) != 3 {
		t.Fatal("clone generic message(struct) with slice failed, object length not correct")
	}

	// ############# test array #############
	gmList = &comp.GenericMessage{
		Subtype: "test_clone_array",
		Obj:     testPtrArray,
	}
	cgm = gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message with slice failed")
	}
	listObj = cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message with slice failed, element not cloned")
	}
	ctlSlice, ok := listObj.([3]*testCloneableObjPtr)
	if !ok {
		t.Fatal("clone generic message with slice failed, must be array type")
	}
	if len(ctlSlice) != 3 {
		t.Fatal("clone generic message with slice failed, element length not correct")
	}

}

func TestGenericMessage_Primitives(t *testing.T) {
	// ############ test primitives ###########
	primList := []interface{}{
		int(1), uint(1), "somestring",
		int8(1), int16(1), int32(1), int64(1),
		uint8(1), uint16(1), uint32(1), uint64(1),
		float32(1.0), float64(1.0),
	}
	for _, p := range primList {
		gm := &comp.GenericMessage{
			Subtype: "primitive",
			Obj:     p,
		}
		cgm := gm.Clone()
		if cgm == nil {
			t.Fatalf("failed to clone primitive object with type: %v", reflect.TypeOf(p))
		}

		gm = &comp.GenericMessage{
			Subtype: "primitive_ptr",
			Obj:     &p,
		}

		cgm = gm.Clone()
		if cgm == nil {
			t.Fatalf("failed to clone primitive object ptr with type: %v", reflect.TypeOf(&p))
		}

		//newObj := cgm.(*comp.GenericMessage).Obj
		//t.Logf("type: %v ---> %v\n", reflect.ValueOf(&p).Elem().Elem().Type(),
		//	reflect.ValueOf(newObj).Elem().Elem().Type())

	}
}

func TestGenericMessage_ConvertibleObject(t *testing.T) {
	obj := &testCloneableObjPtr{1}
	gm := &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     obj,
	}
	cgm := gm.Clone()
	if cgm == nil {
		t.Fatal("generic message does not clone its internal object")
	}
	newObj := cgm.(*comp.GenericMessage).Obj
	if _, ok := newObj.(*testCloneableObjPtr); !ok {
		t.Fatalf("failed to convert to original type, %v", newObj)
	}

	// test slice type with nil pointer
	objSlice := []*testCloneableObjPtr{
		&testCloneableObjPtr{1},
		nil,
		&testCloneableObjPtr{2},
	}
	gm = &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     objSlice,
	}
	cgm = gm.Clone()
	if cgm == nil {
		t.Fatal("generic message does not clone its internal object (array)")
	}
	newSliceObj := cgm.(*comp.GenericMessage).Obj
	if _, ok := newSliceObj.([]*testCloneableObjPtr); !ok {
		t.Fatalf("failed to convert to original slice of type, %v", newSliceObj)
	}

	// test array type with some nil pointers
	objArray := [4]*testCloneableObjPtr{
		&testCloneableObjPtr{1},
		nil,
		nil,
		&testCloneableObjPtr{2},
	}
	gm = &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     objArray,
	}
	cgm = gm.Clone()
	if cgm == nil {
		t.Fatal("generic message does not clone its internal object (array)")
	}
	newArrayObj := cgm.(*comp.GenericMessage).Obj
	if _, ok := newArrayObj.([4]*testCloneableObjPtr); !ok {
		t.Fatalf("failed to convert to original array of type, %v", newSliceObj)
	}
}
