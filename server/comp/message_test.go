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
	return &testCloneableObjPtr{}
}

func (t testCloneableObjStruct) Clone() comp.Cloneable {
	return testCloneableObjStruct{}
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
	testPtrList := []*testCloneableObjPtr{
		&testCloneableObjPtr{1},
		&testCloneableObjPtr{2},
		&testCloneableObjPtr{3},
	}

	testStructList := []testCloneableObjStruct{
		{1},
		{2},
		{3},
	}

	// ############# test array #############
	gmList := &comp.GenericMessage{
		Subtype: "test_clone_array_ptr",
		Obj:     testPtrList,
	}
	cgm := gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message with list failed")
	}
	listObj := cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message with list failed, element not cloned")
	}
	ctl, ok := listObj.([]interface{})
	if !ok {
		t.Fatal("clone generic message with list failed, must be array type")
	}
	if len(ctl) != 3 {
		t.Fatal("clone generic message with list failed, element length not correct")
	}
	if _, ok = ctl[0].(*testCloneableObjPtr); !ok {
		t.Fatalf("clone generic message with list failed, element type not correct, which actually is %v",
			reflect.TypeOf(ctl[0]))
	}

	gmList = &comp.GenericMessage{
		Subtype: "test_clone_array_struct",
		Obj:     testStructList,
	}
	cgm = gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message(struct) with list failed")
	}
	listObj = cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message(struct) with list failed, element not cloned")
	}
	ctl, ok = listObj.([]interface{})
	if !ok {
		t.Fatal("clone generic message(struct) with list failed, must be array type")
	}
	if len(ctl) != 3 {
		t.Fatal("clone generic message(struct) with list failed, element length not correct")
	}
	if _, ok = ctl[0].(testCloneableObjStruct); !ok {
		t.Fatalf("clone generic message(struct) with list failed, element type not correct, which actually is %v",
			reflect.TypeOf(ctl[0]))
	}

	// ############# test slice #############
	gmList = &comp.GenericMessage{
		Subtype: "test_clone_array",
		Obj:     testPtrList[:2],
	}
	cgm = gmList.Clone()
	if cgm == nil {
		t.Fatal("clone generic message with slice failed")
	}
	listObj = cgm.(*comp.GenericMessage).Obj
	if listObj == nil {
		t.Fatal("clone generic message with slice failed, element not cloned")
	}
	ctl, ok = listObj.([]interface{})
	if !ok {
		t.Fatal("clone generic message with slice failed, must be array type")
	}
	if len(ctl) != 2 {
		t.Fatal("clone generic message with slice failed, element length not correct")
	}
	if _, ok = ctl[0].(*testCloneableObjPtr); !ok {
		t.Fatal("clone generic message with slice failed, element type not correct")
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
